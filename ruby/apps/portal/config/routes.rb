Rails.application.routes.draw do
  root 'dashboards#show'
  resources :patents, only: :index
end
