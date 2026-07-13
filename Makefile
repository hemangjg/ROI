.PHONY: dev dev-native dev-core dev-full down logs init-local ws6 ws7 ws8 ws9 ws10 ws11 ws12 p2-ws1 p2-ws2 p2-ws3 p2-ws4 p2-ws5 p2-ws6 p2-ws7 p2-ws8 p2-ws9 certs ci

dev:
	docker compose --profile services up --build

dev-native:
	bash scripts/dev.sh

ci:
	bash scripts/ci-local.sh

dev-core:
	docker compose --profile core up -d

dev-full:
	docker compose --profile full up --build

down:
	docker compose --profile full down --remove-orphans

logs:
	docker compose logs -f --tail=100

init-local:
	bash scripts/init-local.sh

ws6:
	bash scripts/setup-ws6.sh

ws7:
	bash scripts/setup-ws7.sh

ws8:
	bash scripts/setup-ws8.sh

ws9:
	bash scripts/setup-ws9.sh

ws10:
	bash scripts/setup-ws10.sh

ws11:
	bash scripts/setup-ws11.sh

ws12:
	bash scripts/setup-ws12.sh

p2-ws1:
	bash scripts/setup-p2-ws1.sh

p2-ws2:
	bash scripts/setup-p2-ws2.sh

p2-ws3:
	bash scripts/setup-p2-ws3.sh

p2-ws4:
	bash scripts/setup-p2-ws4.sh

p2-ws5:
	bash scripts/setup-p2-ws5.sh

p2-ws6:
	bash scripts/setup-p2-ws6.sh

p2-ws7:
	bash scripts/setup-p2-ws7.sh

p2-ws8:
	bash scripts/setup-p2-ws8.sh

p2-ws9:
	bash scripts/setup-p2-ws9.sh

certs:
	bash scripts/gen-local-certs.sh