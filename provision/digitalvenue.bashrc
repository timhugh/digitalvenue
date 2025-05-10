# add current bin to path
PATH=/opt/digitalvenue/current/bin:$PATH

# source current env
export $(cat /opt/digitalvenue/current/env | xargs)
