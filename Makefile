DELIVERY_PKG = internal/pkg/delivery
DOMAIN_PKG = internal/pkg/domain
API_URL = http://localhost:5000/api

.PHONY: easyjson fill func perf clear run

default:

build:
	docker build -t park .

run: build
	docker run --rm -p 5000:5000 --name dbms park

easyjson:
	easyjson -gen_build_flags="-mod=mod" -lower_camel_case -no_std_marshalers -pkg ${DELIVERY_PKG} ${DOMAIN_PKG}

func:
	./technopark-dbms-forum func -u ${API_URL} -k -r report.html

fill:
	./technopark-dbms-forum fill -u ${API_URL} --timeout=900

perf: fill
	./technopark-dbms-forum perf -u ${API_URL} --duration=600 --step=60

clear:
	curl -X POST ${API_URL}/service/clear
