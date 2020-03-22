class PatentsController < ApplicationController
  def index
    svc = PatentService.new
    @patents = svc.list_patents.patents
  end
end
