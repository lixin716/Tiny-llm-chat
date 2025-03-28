# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import llm_service_pb2 as llm__service__pb2


class LLMServiceStub(object):
    """然后定义服务，使用不同的方法名
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Generate = channel.unary_unary(
                '/llm.LLMService/Generate',
                request_serializer=llm__service__pb2.GenerateRequest.SerializeToString,
                response_deserializer=llm__service__pb2.GenerateResponse.FromString,
                )


class LLMServiceServicer(object):
    """然后定义服务，使用不同的方法名
    """

    def Generate(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_LLMServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Generate': grpc.unary_unary_rpc_method_handler(
                    servicer.Generate,
                    request_deserializer=llm__service__pb2.GenerateRequest.FromString,
                    response_serializer=llm__service__pb2.GenerateResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'llm.LLMService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class LLMService(object):
    """然后定义服务，使用不同的方法名
    """

    @staticmethod
    def Generate(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/llm.LLMService/Generate',
            llm__service__pb2.GenerateRequest.SerializeToString,
            llm__service__pb2.GenerateResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
