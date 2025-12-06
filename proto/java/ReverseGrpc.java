import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * =============================
 * Gufo gRPC Service Definition
 * =============================
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler",
    comments = "Source: microservice.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class ReverseGrpc {

  private ReverseGrpc() {}

  public static final String SERVICE_NAME = "Reverse";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<Microservice.Request,
      Microservice.Response> getDoMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Do",
      requestType = Microservice.Request.class,
      responseType = Microservice.Response.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<Microservice.Request,
      Microservice.Response> getDoMethod() {
    io.grpc.MethodDescriptor<Microservice.Request, Microservice.Response> getDoMethod;
    if ((getDoMethod = ReverseGrpc.getDoMethod) == null) {
      synchronized (ReverseGrpc.class) {
        if ((getDoMethod = ReverseGrpc.getDoMethod) == null) {
          ReverseGrpc.getDoMethod = getDoMethod =
              io.grpc.MethodDescriptor.<Microservice.Request, Microservice.Response>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Do"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  Microservice.Request.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  Microservice.Response.getDefaultInstance()))
              .setSchemaDescriptor(new ReverseMethodDescriptorSupplier("Do"))
              .build();
        }
      }
    }
    return getDoMethod;
  }

  private static volatile io.grpc.MethodDescriptor<Microservice.Request,
      Microservice.Response> getStreamMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Stream",
      requestType = Microservice.Request.class,
      responseType = Microservice.Response.class,
      methodType = io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
  public static io.grpc.MethodDescriptor<Microservice.Request,
      Microservice.Response> getStreamMethod() {
    io.grpc.MethodDescriptor<Microservice.Request, Microservice.Response> getStreamMethod;
    if ((getStreamMethod = ReverseGrpc.getStreamMethod) == null) {
      synchronized (ReverseGrpc.class) {
        if ((getStreamMethod = ReverseGrpc.getStreamMethod) == null) {
          ReverseGrpc.getStreamMethod = getStreamMethod =
              io.grpc.MethodDescriptor.<Microservice.Request, Microservice.Response>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Stream"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  Microservice.Request.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  Microservice.Response.getDefaultInstance()))
              .setSchemaDescriptor(new ReverseMethodDescriptorSupplier("Stream"))
              .build();
        }
      }
    }
    return getStreamMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static ReverseStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ReverseStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ReverseStub>() {
        @java.lang.Override
        public ReverseStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ReverseStub(channel, callOptions);
        }
      };
    return ReverseStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static ReverseBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ReverseBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ReverseBlockingStub>() {
        @java.lang.Override
        public ReverseBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ReverseBlockingStub(channel, callOptions);
        }
      };
    return ReverseBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static ReverseFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ReverseFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ReverseFutureStub>() {
        @java.lang.Override
        public ReverseFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ReverseFutureStub(channel, callOptions);
        }
      };
    return ReverseFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * =============================
   * Gufo gRPC Service Definition
   * =============================
   * </pre>
   */
  public static abstract class ReverseImplBase implements io.grpc.BindableService {

    /**
     */
    public void do_(Microservice.Request request,
        io.grpc.stub.StreamObserver<Microservice.Response> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getDoMethod(), responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<Microservice.Request> stream(
        io.grpc.stub.StreamObserver<Microservice.Response> responseObserver) {
      return io.grpc.stub.ServerCalls.asyncUnimplementedStreamingCall(getStreamMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getDoMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                Microservice.Request,
                Microservice.Response>(
                  this, METHODID_DO)))
          .addMethod(
            getStreamMethod(),
            io.grpc.stub.ServerCalls.asyncBidiStreamingCall(
              new MethodHandlers<
                Microservice.Request,
                Microservice.Response>(
                  this, METHODID_STREAM)))
          .build();
    }
  }

  /**
   * <pre>
   * =============================
   * Gufo gRPC Service Definition
   * =============================
   * </pre>
   */
  public static final class ReverseStub extends io.grpc.stub.AbstractAsyncStub<ReverseStub> {
    private ReverseStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ReverseStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ReverseStub(channel, callOptions);
    }

    /**
     */
    public void do_(Microservice.Request request,
        io.grpc.stub.StreamObserver<Microservice.Response> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getDoMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<Microservice.Request> stream(
        io.grpc.stub.StreamObserver<Microservice.Response> responseObserver) {
      return io.grpc.stub.ClientCalls.asyncBidiStreamingCall(
          getChannel().newCall(getStreamMethod(), getCallOptions()), responseObserver);
    }
  }

  /**
   * <pre>
   * =============================
   * Gufo gRPC Service Definition
   * =============================
   * </pre>
   */
  public static final class ReverseBlockingStub extends io.grpc.stub.AbstractBlockingStub<ReverseBlockingStub> {
    private ReverseBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ReverseBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ReverseBlockingStub(channel, callOptions);
    }

    /**
     */
    public Microservice.Response do_(Microservice.Request request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getDoMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * =============================
   * Gufo gRPC Service Definition
   * =============================
   * </pre>
   */
  public static final class ReverseFutureStub extends io.grpc.stub.AbstractFutureStub<ReverseFutureStub> {
    private ReverseFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ReverseFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ReverseFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<Microservice.Response> do_(
        Microservice.Request request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getDoMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_DO = 0;
  private static final int METHODID_STREAM = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final ReverseImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(ReverseImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_DO:
          serviceImpl.do_((Microservice.Request) request,
              (io.grpc.stub.StreamObserver<Microservice.Response>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_STREAM:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.stream(
              (io.grpc.stub.StreamObserver<Microservice.Response>) responseObserver);
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class ReverseBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    ReverseBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return Microservice.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Reverse");
    }
  }

  private static final class ReverseFileDescriptorSupplier
      extends ReverseBaseDescriptorSupplier {
    ReverseFileDescriptorSupplier() {}
  }

  private static final class ReverseMethodDescriptorSupplier
      extends ReverseBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    ReverseMethodDescriptorSupplier(String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (ReverseGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new ReverseFileDescriptorSupplier())
              .addMethod(getDoMethod())
              .addMethod(getStreamMethod())
              .build();
        }
      }
    }
    return result;
  }
}
