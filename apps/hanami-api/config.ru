require 'bundler/setup'
require 'hanami/api'

class App < Hanami::API
  get '/' do
    ""
  end
end

run App.new
