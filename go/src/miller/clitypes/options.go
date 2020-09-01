package clitypes

// Items which might better belong in miller/cli, but which are placed in a
// deeper package to avoid a package-dependency cycle between miller/cli and
// miiller/mapping.


// ----------------------------------------------------------------
//typedef struct _generator_opts_t {
//	char* field_name;
//	// xxx to do: convert to mv_t
//	long long start;
//	long long stop;
//	long long step;
//} generator_opts_t;

// ----------------------------------------------------------------
type TReaderOptions struct {
	InputFileFormat string
	IRS             string
	IFS             string
	IPS             string

	//	char* input_json_flatten_separator;
	//	json_array_ingest_t  json_array_ingest;

	//	allow_repeat_ifs bool;
	//	allow_repeat_ips bool;
	//	use_implicit_csv_header bool;
	//	allow_ragged_csv_input bool;
	//
	//	// Command for popen on input, e.g. "zcat -cf <". Can be null in which case
	//	// files are read directly rather than through a pipe.
	//	prepipe string;
	//
	//	comment_handling_t comment_handling;
	//	comment_string string;
	//
	//	// Fake internal-data-generator 'reader'
	//	generator_opts_t generator_opts;
}

// ----------------------------------------------------------------
type TWriterOptions struct {
	OutputFileFormat string
	ORS              string
	OFS              string
	OPS              string

	//	headerless_csv_output bool;
	//	right_justify_xtab_value bool;
	//	right_align_pprint bool;
	//	pprint_barred bool;
	//	stack_json_output_vertically bool;
	//	wrap_json_output_in_outer_list bool;
	//	json_quote_int_keys bool;
	//	json_quote_non_string_values bool;
	//	output_json_flatten_separator string;
	//	oosvar_flatten_separator string;
	//
	//	quoting_t oquoting;
}

// ----------------------------------------------------------------
type TOptions struct {
	ReaderOptions TReaderOptions
	WriterOptions TWriterOptions

	// These are used to construct the mapper list. In particular, for in-place mode
	// they're reconstructed for each file.  We make copies since each pass through a
	// CLI-parser operates destructively, principally by running strtok over
	// comma-delimited field-name lists.
	//
	//	char**  original_argv;
	//	char**  non_in_place_argv;
	//	int     argc;
	//	int     mapper_argb;
	//
	//	filenames []string;
	//
	//	char* ofmt;
	//	nr_progress_mod int64u;
	//
	//	do_in_place bool;
	//
	//	no_input bool;
	//	have_rand_seed bool;
	//	rand_seed uint32;
}

// ----------------------------------------------------------------
func DefaultOptions() TOptions {
	return TOptions{
		ReaderOptions: DefaultReaderOptions(),
		WriterOptions: DefaultWriterOptions(),
	}
}

func DefaultReaderOptions() TReaderOptions {
	return TReaderOptions{
		InputFileFormat: "dkvp", // xxx constify at top
		IRS:             "\n",
		IFS:             ",",
		IPS:             "=",
	}
}

func DefaultWriterOptions() TWriterOptions {
	return TWriterOptions{
		OutputFileFormat: "dkvp",
		ORS:              "\n",
		OFS:              ",",
		OPS:              "=",
	}
}