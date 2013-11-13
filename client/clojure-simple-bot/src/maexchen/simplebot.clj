(ns maexchen.simplebot
  (:require [maexchen.simplebot.system :as system]))

; start using "lein run localhost 9000 john"
(defn -main [hostname port botname]
  (-> (system/make-system hostname (Integer/parseInt port) botname)
      (system/register)
      (system/start)))
