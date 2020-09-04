require 'roda'

class App < Roda
  route do |r|
    r.root do
      ''
    end
  end
end

run App.app
