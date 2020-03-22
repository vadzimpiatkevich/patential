require 'proto/svc/patent/service_services_pb'

class PatentService
  def initialize(host: ENV.fetch('PATENT_SVC_HOST'))
    @svc = Patent::Service::Stub.new(
      host.to_s, :this_channel_is_insecure
    )
  end

  def list_patents
    Rails.logger.info('Requesting patent svc to list patents')
    req = Patent::ListPatentsRequest.new
    svc.list_patents(req)
  end

  private

  attr_reader :svc
end
