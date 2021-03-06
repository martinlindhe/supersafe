# About

supersafe is a small server to be run inside a VM or otherwise controlled
environment.

It serves as a target for running a supplied application,
returning the stdout of the execution.

This is all served over http with no auth, for extra safety.

The intended usage is for probing a Windows VM about the execution of
a given code snippet.

This is currently used by [dustbox-rs/fuzzer](https://github.com/martinlindhe/dustbox-rs/tree/master/fuzzer).

## Usage

The intended target is a 32-bit Windows environment.

The resulting executable runs in Windows XP using go 1.10, which is the
last version that produces Windows programs compatible with Windows XP.

To compile the executable, run

```bash
make windows-exe
````

## The /run endpoint

Executes uploaded program and return stdout

## Example

```bash
curl -F "com=@utils/prober/prober.com" http://supersafe-host:28111/run
```

```rust
// in rust
fn stdout_from_vm_http(prober_com: &str) -> String {
    use curl::easy::{Easy, Form};
    let mut dst = Vec::new();
    let mut easy = Easy::new();
    easy.url("http://10.10.30.63:28111/run").unwrap();

    let mut form = Form::new();
    form.part("com").file(prober_com).add().unwrap();
    easy.httppost(form).unwrap();

    {
        let mut transfer = easy.transfer();
        transfer.write_function(|data| {
            dst.extend_from_slice(data);
            Ok(data.len())
        }).unwrap();
        transfer.perform().unwrap();
    }

    str::from_utf8(&dst).unwrap().to_owned()
}
```
