(ns maexchen.simplebot.bot
  (:require [clojure.string :as string]))

(defn- token [parts]
  (nth parts 2))

(defn send-msg [bot & msgparts]
  (.send (:sender bot) (string/join ";" msgparts))
  bot)

(defn round-starting [bot parts]
  (send-msg bot "JOIN" (second parts)))

(defn turn [this parts]
  (send-msg this "ROLL" (second parts)))

(defn announce [this announcement token]
  (send-msg this "ANNOUNCE" announcement token))

(defn rolled [bot parts]
  (announce bot (second parts) (token parts)))

(defn on-message [bot msg]
  (println "server:" msg)
  (let [parts (string/split msg #";" )]
    (condp = (first parts)
      "ROUND STARTING"    (round-starting bot parts)
      "YOUR TURN"         (turn bot parts)
      "ROLLED"            (rolled bot parts)
      bot)))

(defrecord Bot [botname communicator sender])

(defn make-bot [botname communicator sender]
  (->Bot botname communicator sender))

(defn register [bot]
  (send-msg bot "REGISTER" (:botname bot))
  bot)

