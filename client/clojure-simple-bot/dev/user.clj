(ns user
  "Tools for interactive development with the REPL. This file should
   not be included in a production build of the application."
  (:require
    [clojure.java.io :as io]
    [clojure.java.javadoc :refer (javadoc)]
    [clojure.pprint :refer (pprint)]
    [clojure.reflect :refer (reflect)]
    [clojure.repl :refer (apropos dir doc find-doc pst source)]
    [clojure.set :as set]
    [clojure.string :as str]
    [clojure.test :as test]
    [clojure.tools.namespace.repl :refer (refresh refresh-all)]
    [maexchen.simplebot.system :as system])
  (:import [udphelper UdpCommunicator MessageListener MessageSender]
           [java.lang Thread]))

(def system
  "A Var containing an object representing the application under
   development."
  nil)

(def mainloop nil)

(defn init
  "Creates and initializes the system under development in the Var
   #'system."
  []
  (alter-var-root #'system (constantly (system/make-system "localhost" 9000 "clj-simplebot"))))

(defn start
  "Starts the system running, updates the Var #'system."
  []
  (alter-var-root #'system system/register)
  (alter-var-root #'mainloop (constantly (new Thread #(system/start system))))
  ; Start the mainloop in the background so we can still work in the repl
  (.start mainloop))

(defn stop
  "Stops the system if it is currently running, updates the Var
   #'system."
  []
  (.stop mainloop)
  (alter-var-root #'mainloop (constantly nil))
  (alter-var-root #'system system/stop))

(defn go
  "Initializes and starts the system running."
  []
  (if (nil? mainloop) 
    (do
      (init)
      (start)
      :ready)
    (println "WARNING mainloop running.")))

(defn reset
  "Stops the system, reloads modified source files, and restarts it."
  []
  (stop)
  (refresh :after 'user/go))

