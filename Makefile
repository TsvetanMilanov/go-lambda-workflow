.PHONY: \
	test

test:
	@go test -v -cover github.com/TsvetanMilanov/go-lambda-workflow/workflow
