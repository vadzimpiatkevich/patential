# Generated by the protocol buffer compiler. DO NOT EDIT!
# Source: golang/svc/patent/proto/service/service.proto.

require 'grpc'
require_relative 'service_pb'

module Patent
  module Service
    class Service

      include GRPC::GenericService

      self.marshal_class_method = :encode
      self.unmarshal_class_method = :decode
      self.service_name = 'patent.Service'

      rpc :ListPatents, ListPatentsRequest, ListPatentsResponse
    end

    Stub = Service.rpc_stub_class
  end
end
