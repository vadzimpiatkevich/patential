class PatentsController < ApplicationController
  before_action :authenticate_user!

  def index
    svc = PatentService.new
    @patents = svc.list_patents.patents
  end
end
