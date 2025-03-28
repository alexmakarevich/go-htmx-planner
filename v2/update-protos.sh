rm -f proto-go/*;
# rm -f svelte/proto-ts/*
rm -f svelte/proto-es/*

protoc --proto_path=proto --go_out=proto-go --go_opt=paths=source_relative proto/*.proto;
# protoc  --plugin=./svelte/node_modules/.bin/protoc-gen-ts_proto --proto_path=proto --ts_proto_out=./svelte/proto-ts proto/*.proto;
cd svelte && npx buf generate;


# TODO: remove ts-proto, if es-proto is enough