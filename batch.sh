#!/bin/bash

#!/bin/bash
for i in {1..255}
do
   echo "iteration $i"
   curl -F "com=@/Users/m/dev/rs/dustbox-rs/utils/prober/prober.com" http://10.10.30.63:28111/run
done
