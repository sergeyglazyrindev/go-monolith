TEST_PATTERN?=
RACE?=-race
CURRENT_DIRECTORY?=$(shell pwd)
TEST_ENVIRONMENT?=test.sqlite
UT_PACKAGES?=$(shell $(GO) list ./...)
FUNC_TESTS_CMD:="grep -e 'func Test${TEST_PATTERN}' tests/*.go | perl -pe 's|.*func (.*?)\(.*|\1|g' | shuf"
FUNC_TESTS:=$(shell sh -c $(FUNC_TESTS_CMD))
VERBOSE_TESTS_FLAGS?=
TIMEOUT?=1m

ifeq ($(VERBOSE), true)
  VERBOSE_TESTS_FLAGS+=-test.v
endif

ifeq ($(COVERAGE), true)
  TEST_COVERPROFILE?=../functionals.cover
  EXTRA_ARGS+=-test.coverprofile=${TEST_COVERPROFILE}
endif

COVERAGE?=0
COVERAGE_MODE?=atomic
COVERAGE_WD?="."

comma:= ,
empty:=
space:= $(empty) $(empty)

.PHONY: test.functionals.clean
test.functionals.clean:
	rm -f tests/functionals

.PHONY: test.functionals.compile
test.functionals.compile: genlocalfiles
	$(GO) test -tags "${BUILD_TAGS} test" -race ${GOFLAGS} ${VERBOSE_FLAGS} -timeout ${TIMEOUT} -c -o tests/functionals ./tests/

.PHONY: test.functionals.static
test.functionals.static: genlocalfiles

ifeq (${DEBUG}, true)
define functionals_run
cd tests && sudo -E $$(which dlv) $(DLV_FLAGS) exec ./functionals -- $1
endef
else
define functionals_run
cd tests && sudo -E ./functionals $1
endef
endif

.PHONY: test.functionals.run
test.functionals.run:
	cd tests && sudo -E ./functionals ${VERBOSE_TESTS_FLAGS} -test.run "${TEST_PATTERN}" -test.timeout ${TIMEOUT} ${ARGS} ${EXTRA_ARGS}

.PHONY: tests.functionals.all
test.functionals.all: test.functionals.compile
	$(MAKE) TIMEOUT="8m" ARGS="${ARGS}" test.functionals.run EXTRA_ARGS="${EXTRA_ARGS}"

.PHONY: test.functionals.batch
test.functionals.batch: test.functionals.compile
	set -e ; $(MAKE) ARGS="${ARGS} " test.functionals.run EXTRA_ARGS="${EXTRA_ARGS}" TEST_PATTERN="${TEST_PATTERN}"

.PHONY: test.functionals
test.functionals: test.functionals.compile
	for functest in ${FUNC_TESTS} ; do \
		$(MAKE) ARGS="-test.run $$functest$$\$$ ${ARGS}" test.functionals.run EXTRA_ARGS="${EXTRA_ARGS}"; \
	done

.PHONY: functional
functional:
	$(MAKE) test.functionals VERBOSE=true TIMEOUT=10m ARGS='-standalone' TEST_PATTERN="${TEST_PATTERN}"

.PHONY: test
test: genlocalfiles
ifeq ($(COVERAGE), true)
	set -v ; \
	for pkg in ${UT_PACKAGES}; do \
		if [ -n "$$pkg" ]; then \
			coverfile="${COVERAGE_WD}/go-monolith.cover"; \
        	export GOMONOLITH_PATH=${CURRENT_DIRECTORY} ; \
			export TEST_ENVIRONMENT=${TEST_ENVIRONMENT} ; \
			$(GO) test -tags "${BUILD_TAGS} test" -coverpkg ./... -coverprofile="$$coverfile" ${VERBOSE_FLAGS} -timeout ${TIMEOUT} $$pkg || exit 1; \
		fi; \
	done
else
ifneq ($(TEST_PATTERN),)
	set -v ; \
	export GOMONOLITH_PATH=${CURRENT_DIRECTORY} ; \
	export TEST_ENVIRONMENT=${TEST_ENVIRONMENT} ; \
	$(GO) clean -testcache ; \

	# $(GO) test -tags "${BUILD_TAGS} test" -ldflags="${LDFLAGS}" ${RACE} ${GOFLAGS} ${VERBOSE_FLAGS} -timeout ${TIMEOUT} -test.run ${TEST_PATTERN} ${UT_PACKAGES}
	for pkg in ${UT_PACKAGES}; do \
		if [ -n "$$pkg" ]; then \
			$(GO) test -tags "${BUILD_TAGS} test" -ldflags="${LDFLAGS}" ${RACE} ${GOFLAGS} ${VERBOSE_FLAGS} -timeout ${TIMEOUT} -test.run ${TEST_PATTERN} $$pkg || exit 1; \
		fi; \
	done
else
	set -v ; \
	export GOMONOLITH_PATH=${CURRENT_DIRECTORY} ; \
	export TEST_ENVIRONMENT=${TEST_ENVIRONMENT} ; \
	$(GO) clean -testcache ; \
	# $(GO) test -tags "${BUILD_TAGS} test" -ldflags="${LDFLAGS}" ${RACE} ${GOFLAGS} ${VERBOSE_FLAGS} -timeout ${TIMEOUT} ${UT_PACKAGES}
	for pkg in ${UT_PACKAGES}; do \
		if [ -n "$$pkg" ]; then \
			$(GO) test -tags "${BUILD_TAGS} test" -ldflags="${LDFLAGS}" ${RACE} ${GOFLAGS} ${VERBOSE_FLAGS} -timeout ${TIMEOUT} $$pkg || exit 1; \
		fi; \
	done
endif
endif
