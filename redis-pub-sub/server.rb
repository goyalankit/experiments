# A simple api to publish message through redis
# sample use:
# http://localhost:4567/publish.json?channel=tm&message=hello
#
# To run it as a daemon
# nohup ruby server.rb >> /var/log/server.log 2>&1 &
#

require 'sinatra'
require 'redis'
require 'json'

set :bind, '0.0.0.0'
$redis = Redis.new


get '/publish.json' do
  # we publish all params
  $redis.publish params["channel"], params.to_json
end
