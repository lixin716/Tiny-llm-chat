# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: llm_service.proto
# Protobuf Python Version: 4.25.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x11llm_service.proto\x12\x03llm\"]\n\x0fGenerateRequest\x12\x0e\n\x06prompt\x18\x01 \x01(\t\x12\x13\n\x0btemperature\x18\x02 \x01(\x02\x12\x16\n\x0emax_new_tokens\x18\x03 \x01(\x05\x12\r\n\x05top_k\x18\x04 \x01(\x05\"$\n\x10GenerateResponse\x12\x10\n\x08response\x18\x01 \x01(\t2G\n\nLLMService\x12\x39\n\x08Generate\x12\x14.llm.GenerateRequest\x1a\x15.llm.GenerateResponse\"\x00\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'llm_service_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:
  DESCRIPTOR._options = None
  _globals['_GENERATEREQUEST']._serialized_start=26
  _globals['_GENERATEREQUEST']._serialized_end=119
  _globals['_GENERATERESPONSE']._serialized_start=121
  _globals['_GENERATERESPONSE']._serialized_end=157
  _globals['_LLMSERVICE']._serialized_start=159
  _globals['_LLMSERVICE']._serialized_end=230
# @@protoc_insertion_point(module_scope)
