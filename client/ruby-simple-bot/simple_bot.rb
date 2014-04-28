require_relative '../ruby-udp-helper/udp_connector'

class SimpleBot
  def initialize(connection)
    @connection = connection
  end

  def register name
    @connection.send_message('REGISTER;' + name)
    response = @connection.receive_message

    if response == 'REGISTERED'
      play
    else
      puts 'ERROR: Unable to connect!'
    end
  end

  def play
    while true
      response = @connection.receive_message
      make_move response
    end
  end

  def make_move response
    parts = response.split(';')

    case parts[0]
    when 'ROUND STARTING'
      @connection.send_message('JOIN;' + parts[1])
    when 'YOUR TURN'
      @connection.send_message('ROLL;' + parts[1])
    when 'ROLLED'
      @connection.send_message('ANNOUNCE;' + parts[1] + ';' + parts[2])
    end
  end
end

if ARGV.length == 3
  host, port, name = ARGV

  client = UdpConnector.new(host, port)
  bot = SimpleBot.new(client)
  bot.register name
else
  puts 'USAGE: ruby simple_bot.rb HOST PORT NAME'
end
