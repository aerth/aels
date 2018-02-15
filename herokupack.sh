#!/bin/bash
##################
# aerth waz here #
##################
set +e
pack_vendor(){
	echo "removing 'vendor' directory"	
	rm -rf vendor
	echo "downloading vendor repos... may take a few minutes."
	printf "fetching external packages..."
	govendor init && govendor fetch +external
	printf "all done!\n"
}

check_Procfile(){
	test -f Procfile && echo "Procfile found."
	pf=$(cat Procfile)
	printf "Contents of Procfile:\n$pf\n"
}

test -x $(which govendor) || (echo "need govendor installed." && exit 1)
pack_vendor
check_Procfile
