#
# This spec file is read by gdc when linking.
# It is used to specify the libraries we need to link in, in the right
# order.
#


%rename link linkorig_gdc_renamed
*link: %(linkorig_gdc_renamed) 

%rename lib liborig_gdc_renamed
*lib: %{debuglib|defaultlib|nophoboslib: ; : } %(liborig_gdc_renamed)
