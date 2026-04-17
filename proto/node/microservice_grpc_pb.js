// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var microservice_pb = require('./microservice_pb.js');
var google_protobuf_any_pb = require('google-protobuf/google/protobuf/any_pb.js');
var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js');

function serialize_gufo_Request(arg) {
  if (!(arg instanceof microservice_pb.Request)) {
    throw new Error('Expected argument of type gufo.Request');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_gufo_Request(buffer_arg) {
  return microservice_pb.Request.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_gufo_Response(arg) {
  if (!(arg instanceof microservice_pb.Response)) {
    throw new Error('Expected argument of type gufo.Response');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_gufo_Response(buffer_arg) {
  return microservice_pb.Response.deserializeBinary(new Uint8Array(buffer_arg));
}


// ============================================================
// Gufo Core Transport (Dynamic Execution Bus)
// ============================================================
// - Keeps URL structure compatibility:
//   /api/v1/{module}/{param}/{param_id}/{param_idd}
// - Supports HTTP + gRPC
// - Separates routing, query, and body
// ============================================================
//
var ReverseService = exports.ReverseService = {
  do: {
    path: '/gufo.Reverse/Do',
    requestStream: false,
    responseStream: false,
    requestType: microservice_pb.Request,
    responseType: microservice_pb.Response,
    requestSerialize: serialize_gufo_Request,
    requestDeserialize: deserialize_gufo_Request,
    responseSerialize: serialize_gufo_Response,
    responseDeserialize: deserialize_gufo_Response,
  },
  stream: {
    path: '/gufo.Reverse/Stream',
    requestStream: true,
    responseStream: true,
    requestType: microservice_pb.Request,
    responseType: microservice_pb.Response,
    requestSerialize: serialize_gufo_Request,
    requestDeserialize: deserialize_gufo_Request,
    responseSerialize: serialize_gufo_Response,
    responseDeserialize: deserialize_gufo_Response,
  },
};

exports.ReverseClient = grpc.makeGenericClientConstructor(ReverseService);
