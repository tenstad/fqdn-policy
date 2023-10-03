#!/bin/zsh

if [ $# -lt 1 ]; then
  echo "Needs parameter add or remove"
  exit 1
fi

removeAnnotation() {
echo "Removing annotations"
while read NS NETPOL; do
  kubectl annotate netpol $NETPOL -n $NS fqdnnetworkpolicies.networking.gke.io/owned-by- 
done < <(kubectl get netpol -A | grep -E '\-fqdn|\-egress' | awk -F" " '{ print $1 " " $2 }')
}

addAnnotation() {
echo "Adding annotations"
while read NS NETPOL; do
  kubectl annotate netpol $NETPOL -n $NS fqdnnetworkpolicies.networking.gke.io/owned-by=$NETPOL
done < <(kubectl get netpol -A | grep -E '\-fqdn|\-egress' | awk -F" " '{ print $1 " " $2 }')
}

if [ "$1" = "add" ]; then
  addAnnotation
elif [ "$1" = "remove" ]; then
  removeAnnotation
else
  echo "Wrong parameter, use add or remove"
fi
