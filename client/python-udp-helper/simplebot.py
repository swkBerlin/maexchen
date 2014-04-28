from __future__ import print_function
import socket
import argparse

class SimpleBot(object):

  def __init__(self, host, port, name):
    self.host = host
    self.port = port
    self.name = name
    self.sock = None

  def start(self):
    self.connect_to_server()
    response = self.receive_message()

    if response == 'REGISTERED':
      print("connected to {}:{}".format(self.host, self.port))
      self.participate()
    else:
      print("server rejected connection")
    self.sock.close()

  def participate(self):
    while True:
      msg = self.receive_message()
      self.react_on_message(msg)

      
  def react_on_message(self, msg):
    parts = msg.split(";")
    command = parts[0]

    if command == "ROUND STARTING":
      self.send_message("JOIN;{}".format(parts[1]))
    elif command == "YOUR TURN":
      self.send_message("ROLL;{}".format(parts[1]));
    elif command == "ROLLED":
      self.send_message("ANNOUNCE;{};{}".format(parts[1], parts[2]));
  
  def connect_to_server(self):
    self.sock = socket.socket(socket.AF_INET,
                          socket.SOCK_DGRAM)
    self.sock.connect((self.host, self.port))
    self.sock.sendall("REGISTER;{}".format(self.name))

  def receive_message(self):
    try:
      msg = self.sock.recv(1024)
      print("received {}".format(msg))
      return msg
    except socket.error, e:
      print("Failed to receive message: {}".format(e))

  def send_message(self, msg):
    try:
      self.sock.sendall(msg)
      print("sent {}".format(msg))
    except Exception, e:
      print("Failed to send message {}: {}".format(msg, e))

if __name__ == "__main__":
  parser = argparse.ArgumentParser(description='Client to connect to Mia server. (https://www.github.com/swkBerlin/maexchen')
  parser.add_argument('host', type=str,
                   help='mia host IP')
  parser.add_argument('port', type=int,
                   help='mia host IP')
  parser.add_argument('name', type=str,
                   help='mia host IP')
  args = parser.parse_args()
  
  bot = SimpleBot(args.host, args.port, args.name)
  bot.start()
