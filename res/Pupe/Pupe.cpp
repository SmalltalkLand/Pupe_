#include <string>
#include <vector>
#include <gsl/gsl>
#include <experimental/coroutine>
#include <gtk/gtk.h>
#include <unistd.h>
#include <gtkmm/button.h>
#include <gtkmm/window.h>
#include <signal.h>
#include "tuple.cpp"
template<typename T,typename... Args>class ManualDestroy{
T *v;
public:
ManualDestroy(Args... args){
    v = new T(args...);
}
void destroy(){
    if(v == 0)return;
    delete v;
    v = 0;
}
}
template<typename F,typename T> T reduce(F fn,T a){
    return a;
}
template<typename F,typename T,typename... Args> T reduce(F fn,T a,Args... args){
    return fn(a,reduce(fn,args...));
}
template<typename F,typename T,typename... Args,typename... Args2> T lambdize(F fn,Args... args){
    return [=](Args2... a){return fn(args...,a...);};
}
template<typename F,typename C,typename T,typename E,typename... Args> T catch_(F fn,C handler,Args... args){
    try{
        return fn(args...);
    }catch(E err){
        return handler(err,args...);
    }
}
namespace pointers{
    class RC
{
    private:
    long long long int count; // Reference count

    public:
    void AddRef()
    {
        // Increment the reference count
        count++;
    }

    int Release()
    {
        // Decrement the reference count and
        // return the reference count.
        return --count;
    }
};
template < typename T > class SP
{
private:
    T*    pData;       // pointer
    RC* reference; // Reference count

public:
    SP() : pData(0), reference(0) 
    {
        // Create a new reference 
        reference = new RC();
        // Increment the reference count
        reference->AddRef();
    }

    SP(T* pValue) : pData(pValue), reference(0)
    {
        // Create a new reference 
        reference = new RC();
        // Increment the reference count
        reference->AddRef();
    }

    SP(const SP<T>& sp) : pData(sp.pData), reference(sp.reference)
    {
        // Copy constructor
        // Copy the data and reference pointer
        // and increment the reference count
        reference->AddRef();
    }

    ~SP()
    {
        // Destructor
        // Decrement the reference count
        // if reference become zero delete the data
        if(reference->Release() == 0)
        {
            delete pData;
            delete reference;
        }
    }

    T& operator* ()
    {
        return *pData;
    }

    T* operator-> ()
    {
        return pData;
    }
    
    SP<T>& operator = (const SP<T>& sp)
    {
        // Assignment operator
        if (this != &sp) // Avoid self assignment
        {
            // Decrement the old reference count
            // if reference become zero delete the old data
            if(reference->Release() == 0)
            {
                delete pData;
                delete reference;
            }

            // Copy the data and reference pointer
            // and increment the reference count
            pData = sp.pData;
            reference = sp.reference;
            reference->AddRef();
        }
        return *this;
    }
};
}
class Proc
{
public:
  Proc( const std::string& cmd, std::vector<std:string> args)
  {
    m_pid = fork()
    if (pid == 0) {
        std::vector<char*> c;
        for (std::string arg: args)c.push_back(arg.c_str());
      execl(cmd.c_str(), &c[0], NULL);
    }
  }
  void kill(int signal){
      kill(m_pid,signal)
      wait(m_pid)
      killed = true
  }
  ~Proc()
  {
    // Just for the case, we have 0, we do not want to kill ourself
    if( m_pid > 0 && !killed)
    {
      kill(m_pid, 9);
      wait(m_pid);
    }
  }
private:
bool killed;
  pid_t m_pid;
}
class PupeBus{
    std::ofstream b;
    public:
    PupeBus(std::string path): b(path){

    }
    void emit(std::string data){
        b << data
    }
}
int m2(std::vector<std::string> argv){
    PupeBus p("/tmp/pupe/bus0");
    std::string self = readlink("/proc/self/exe");

    return 0;
}
int main(int argc,char** argv){
std::vector<char*> a(argv,argv + argc);
std::vector<std::string> b;
for(char* x: a){
    b.push_back(std::string(x));
}
    return m2(b);
}
