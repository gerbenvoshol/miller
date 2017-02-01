#include <stdio.h>
#include <string.h>
#include "lib/context.h"
#include "cli/mlrcli.h"

// ----------------------------------------------------------------
void context_init_from_first_file_name(context_t* pctx, char* first_file_name) {
	memset(pctx, 0, sizeof(context_t));
	pctx->nr       = 0;
	pctx->fnr      = 0;
	pctx->filenum  = 1;
	pctx->filename = first_file_name;
	pctx->auto_line_term = "\n"; // xxx default to "\r\n" on Windows
	pctx->auto_line_term_detected = 0;
}

// ----------------------------------------------------------------
void context_init_from_opts(context_t* pctx, void* pvopts) {
	cli_opts_t* popts = pvopts;

	pctx->nr        = 0;
	pctx->fnr       = 0;
	pctx->filenum   = 0;
	pctx->filename  = NULL;
	pctx->force_eof = 0;

	pctx->ips       = popts->reader_opts.ips;
	pctx->ifs       = popts->reader_opts.ifs;
	pctx->irs       = popts->reader_opts.irs;
	pctx->ops       = popts->writer_opts.ops;
	pctx->ofs       = popts->writer_opts.ofs;
	pctx->ors       = popts->writer_opts.ors;
	pctx->auto_line_term = "\n"; // xxx Windows default "\r\n"; libify
	pctx->auto_line_term_detected = 0;
}

// ----------------------------------------------------------------
void context_print(context_t* pctx, char* indent) {
	printf("%spctx at %p:\n", indent, pctx);
	printf("%s  nr       = %lld\n", indent, pctx->nr);
	printf("%s  fnr      = %lld\n", indent, pctx->fnr);
	printf("%s  filenum  = %d\n", indent, pctx->filenum);
	if (pctx->filename == NULL) {
		printf("%s  filename = null\n", indent);
	} else {
		printf("%s  filename = \"%s\"\n", indent, pctx->filename);
	}
}
