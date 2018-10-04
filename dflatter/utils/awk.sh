awk '!/^FROM/ || !f++' Dockerfile | awk '!/^ARG SSH_AUTH_HOST/ || !f++' | awk '!/^ARG BAY_HOST/ || !f++' | awk '!/MAINTAINER/' | awk  '/^ENTRYPOINT/ { entry = $0; next } { print }  END { print entry }' > Dockerfile.new
rm Dockerfile
mv Dockerfile.new Dockerfile
sed -i 's/ADD /ADD base\//g' Dockerfile
sed -i 's/COPY /COPY base\//g' Dockerfile
sed -i 's/ADD base\/build/ADD build/g' Dockerfile

remove=( "base/Dockerfile-12.04" "base/eb_sources_12.04.list")

for i in "${remove[@]}"
do
   if [ -f $i ] ; then
   echo "removing" $i
    rm $i
   fi
done

cat >> bay.yaml <<_EOF_
boot:
    build:
        apt-cacher: optional
        ssh-agent: required
        squid: optional
    run:
        ssh-agent: required
_EOF_