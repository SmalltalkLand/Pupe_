#include <unistd.h>
#include <string.h>
#include <fcntl.h>
int main2(int argc,char* proc,char** argv,char** env){
    if(strcmp(argv[0],"insert")==0){
        char** argv2 = &argv[1];
        while(strcmp(argv2[0],"--")){
        int h = atoi(argv2[0]);
        int o = open(argv2[1], O_RDONLY);
        dup2(o,h);
        argv2=&argv2[2];
        };
        execve(argv2[0],&argv2[1],&env[0]);
    };
    if(strcmp(argv[0],"recv")==0){
        if(strcmp(argv[1],"-chroot")==0){chroot(argv[2]);return main2(argc - 3,proc,&argv[3],env);};
            char** argv2 = &argv[1];
        while(strcmp(argv2[0],"--")){
        int h = atoi(argv2[0]);
        int o = open(argv2[1], O_WRONLY);
        char buf[256];
        int n;
        while(n=read(h,&buf,sizeof buf)){
            write(o,&buf,n);
            
        };
        argv2=&argv2[2];
        };
        fexecve(atoi(argv2[0]),&argv2[1],&env[0]);
    };
    return 1;
};
int main(int argc,char** argv,char** env){
        char* proc = argv[0];
        return main2(argc,proc,&argv[1],&env[0]);
}
