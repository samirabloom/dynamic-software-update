#!/usr/bin/env sh
 
get_vbox_version(){
    local VER
    VER=$(VBoxManage -v | awk -F "r" '{print $1}')
    if [ -z "$VER" ]; then
        echo "ERROR"
    else
        echo "$VER"
    fi
 
}
 
write_vbox_dockerfile(){
    local VER
    VER=$(get_vbox_version)
    if [ ! "$LATEST_RELEASE" = "ERROR" ]; then
        sed "s/\$VBOX_VERSION/$VER/g" Dockerfile.tmpl > Dockerfile
    else
        echo "WUH WOH"
    fi
}
 
write_vbox_dockerfile