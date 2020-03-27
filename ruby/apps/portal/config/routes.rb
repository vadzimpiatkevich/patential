Rails.application.routes.draw do
  root 'dashboards#show'
  devise_for :users
  resources :patents, only: :index
end
