commons=($(diff -srq ../bundle/core-js ../bundle/core-frontend | grep identical | awk '{print $2}'))

for i in "${commons[@]}"
do
   :
  echo $i
done