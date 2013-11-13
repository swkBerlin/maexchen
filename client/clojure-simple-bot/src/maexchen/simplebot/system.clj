(ns maexchen.simplebot.system
  (:require [maexchen.simplebot.bot :as bot])
  (:import [udphelper UdpCommunicator MessageSender MessageListener]))

; Encapsulate the mutablility in an :bot atom in the system structure
; See bot/on-message for handling server messages

(defn- make-listening-bot [bot]
  (reify MessageListener
    (onMessage [this msg]
      (swap! bot bot/on-message msg))))

(defn make-system [hostname port botname]
  (let [communicator (new UdpCommunicator hostname port)
        sender (.getMessageSender communicator)
        bot (atom (bot/make-bot botname communicator sender))
        listening-bot (make-listening-bot bot)]
    (.addMessageListener communicator listening-bot)
    {
     :communicator communicator
     :sender sender
     :listener listening-bot
     :bot bot
     }))

(defn register [system]
  (swap! (:bot system) bot/register)
  system)

(defn start [system]
  (.listenForMessages (:communicator system)))

(defn stop [system]
  (.removeMessageListener (:communicator system) (:listener system))
  system)

