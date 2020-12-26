#include <string>
#include <vector>
int m2(std::vector<std::string> argv){
    return 0;
}
int main(int argc,char** argv){
std::vector<char*> a(argv,argv + argc);
    return 0;
}