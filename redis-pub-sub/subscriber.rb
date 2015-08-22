# subscriber that listens to publisher
# for a given channel 'tm'.
# change host, port and channel accordingly.

require 'rubygems'
require 'redis'
require 'json'

$redis = Redis.new(:host => "ec2-52-27-241-172.us-west-2.compute.amazonaws.com", :port => 6379)

$redis.subscribe('tm') do |on|
  on.message do |channel, msg|
    require 'pry'; binding.pry
    puts "##{channel} - #{data['message']}"
    # DO SOMETHING CREATIVE HERE
  end
end
