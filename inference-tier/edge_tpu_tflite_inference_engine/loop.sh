#!/bin/bash

# runtime="15000 second"
# endtime=$(date -ud "$runtime" +%s)

# while [[ $(date -u +%s) -le $endtime ]]
# do
#     echo "Time Now: `date +%H:%M:%S`"
#     echo "Sleeping for 1 seconds"
#     # python3 infer_tflite.py
#     sleep 1
#     # python3 infer_tflite_0.py
# done


# echo  "starting py socket server OO..."
python3 infer_tflite_socket.py