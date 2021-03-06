package mappers

import (
	"container/list"
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/mapping"
	"miller/types"
)

// ----------------------------------------------------------------
var RemoveEmptyColumnsSetup = mapping.MapperSetup{
	Verb:         "remove-empty-columns",
	ParseCLIFunc: mapperRemoveEmptyColumnsParseCLI,
	IgnoresInput: false,
}

func mapperRemoveEmptyColumnsParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) mapping.IRecordMapper {

	// Get the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	flagSet := flag.NewFlagSet(verb, errorHandling)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		mapperRemoveEmptyColumnsUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	mapper, _ := NewMapperRemoveEmptyColumns()

	*pargi = argi
	return mapper
}

func mapperRemoveEmptyColumnsUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s {no options}\n", argv0, verb)
	fmt.Fprintf(o, "Omits fields which are empty on every input row. Non-streaming.\n")

	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type MapperRemoveEmptyColumns struct {
	recordsAndContexts      *list.List
	namesWithNonEmptyValues map[string]bool
}

func NewMapperRemoveEmptyColumns() (*MapperRemoveEmptyColumns, error) {
	this := &MapperRemoveEmptyColumns{
		recordsAndContexts:      list.New(),
		namesWithNonEmptyValues: make(map[string]bool),
	}
	return this, nil
}

// ----------------------------------------------------------------
func (this *MapperRemoveEmptyColumns) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		this.recordsAndContexts.PushFront(inrecAndContext)

		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if !pe.Value.IsEmpty() {
				this.namesWithNonEmptyValues[*pe.Key] = true
			}
		}

	} else { // end of record stream

		for e := this.recordsAndContexts.Front(); e != nil; e = e.Next() {
			outrecAndContext := e.Value.(*types.RecordAndContext)
			outrec := outrecAndContext.Record

			newrec := types.NewMlrmapAsRecord()

			for pe := outrec.Head; pe != nil; pe = pe.Next {
				_, ok := this.namesWithNonEmptyValues[*pe.Key]
				if ok {
					// Transferring ownership from old record to new record; no copy needed
					newrec.PutReference(pe.Key, pe.Value)
				}
			}

			outputChannel <- types.NewRecordAndContext(newrec, &outrecAndContext.Context)
		}

		outputChannel <- inrecAndContext // Emit the stream-terminating null record
	}
}
