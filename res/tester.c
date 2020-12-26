int main(int argc,char** argv){
    while(1){
    execv(argv[0],&argv[1]);
    };
}
