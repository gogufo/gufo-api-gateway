// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var microservice_pb = require('./microservice_pb.js');
var google_protobuf_any_pb = require('google-protobuf/google/protobuf/any_pb.js');

function serialize_Request(arg) {
  if (!(arg instanceof microservice_pb.Request)) {
    throw new Error('Expected argument of type Request');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_Request(buffer_arg) {
  return microservice_pb.Request.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_Response(arg) {
  if (!(arg instanceof microservice_pb.Response)) {
    throw new Error('Expected argument of type Response');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_Response(buffer_arg) {
  return microservice_pb.Response.deserializeBinary(new Uint8Array(buffer_arg));
}


// =============================
// Gufo gRPC Service Definition
// =============================
var ReverseService = exports.ReverseService = {
  do: {
    path: '/Reverse/Do',
    requestStream: false,
    responseStream: false,
    requestType: microservice_pb.Request,
    responseType: microservice_pb.Response,
    requestSerialize: serialize_Request,
    requestDeserialize: deserialize_Request,
    responseSerialize: serialize_Response,
    responseDeserialize: deserialize_Response,
  },
  stream: {
    path: '/Reverse/Stream',
    requestStream: true,
    responseStream: true,
    requestType: microservice_pb.Request,
    responseType: microservice_pb.Response,
    requestSerialize: serialize_Request,
    requestDeserialize: deserialize_Request,
    responseSerialize: serialize_Response,
    responseDeserialize: deserialize_Response,
  },
};

exports.ReverseClient = grpc.makeGenericClientConstructor(ReverseService);
