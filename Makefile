PLUGIN_NAME=l2vlan

clean:
	rm -rf ./plugin ./bin
	rm -f ${PLUGIN_NAME}
	docker plugin disable ${PLUGIN_NAME} || true
	docker plugin rm ${PLUGIN_NAME} || true
	docker plugin disable blins1999/${PLUGIN_NAME} || true
	docker plugin rm blins1999/${PLUGIN_NAME} || true
	docker rm -vf tmp || true
	docker rmi ${PLUGIN_NAME}-build-image || true
	docker rmi ${PLUGIN_NAME}:rootfs || true

build:
	docker build -t ${PLUGIN_NAME}-build-image -f Dockerfile.build .
	docker create --name tmp ${PLUGIN_NAME}-build-image
	docker cp tmp:/go/bin/${PLUGIN_NAME} .
	docker rm -vf tmp
	#docker rmi ${PLUGIN_NAME}-build-image
	docker build -t ${PLUGIN_NAME}:rootfs .
	mkdir -p ./plugin/rootfs
	docker create --name tmp ${PLUGIN_NAME}:rootfs
	docker export tmp | tar -x -C ./plugin/rootfs
	cp config.json ./plugin/
	docker rm -vf tmp
	rm -f ${PLUGIN_NAME}

create-plugin:
	docker plugin create blins1999/${PLUGIN_NAME} ./plugin

create-plugin-local:
	docker plugin create ${PLUGIN_NAME} ./plugin

push-plugin:
	docker plugin push blins1999/${PLUGIN_NAME}

rm-plugin:
	docker plugin rm ${PLUGIN_NAME} || true
	docker plugin rm blins1999/${PLUGIN_NAME} || true

push: clean build create-plugin push-plugin rm-plugin clean