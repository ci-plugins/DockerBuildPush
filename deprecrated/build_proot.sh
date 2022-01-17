PROOT_INSTALL_PATH=/data/proot_ins
CUR_PATH=`pwd`

function ins_libtalloc()
{
    #编译列态libtalloc.a库和头文件
    cd ${CUR_PATH}
    rm -rf talloc-2.1.1
    tar -xvf talloc-2.1.1.tar.gz
    cd talloc-2.1.1
    ./configure --prefix="${PROOT_INSTALL_PATH}" && make && make install
    ar qf libtalloc.a bin/default/talloc_3.o
    cp libtalloc.a ${PROOT_INSTALL_PATH}/lib
}

function ins_libarchive()
{
    #静态编译libarchive.a库和头文件
    cd ${CUR_PATH}
    #rm -rf libarchive
    #tar -xvf libarchive-3.5.2.tar.gz
    cd libarchive
    ./build/autogen.sh
    ./configure --prefix="${PROOT_INSTALL_PATH}" --enable-shared=no && make && make install
}

function ins_proot()
{
    #静态链接proot
    cd ${CUR_PATH}
    rm -rf proot-5.1.0
    tar -xvf proot-5.1.0.tar.gz
    cd proot-5.1.0

    LDFLAGS="-L/${PROOT_INSTALL_PATH}/lib -static" PKG_CONFIG_PATH="${PROOT_INSTALL_PATH}/lib/pkgconfig" CFLAGS="-I/${PROOT_INSTALL_PATH}/include" make -C src/ proot GIT=false
    cd ${CUR_PATH}
    cp ./proot-5.1.0/src/proot ./proot
}

function ins_glibc_static()
{
  #yum install -y glibc-devel
  yum install -y python-devel
  yum install -y glibc-static
  yum install -y autoconf automake libtool
}

function download_talloc()
{
    cd ${CUR_PATH}
    curl -o talloc-2.1.1.tar.gz https://www.samba.org/ftp/talloc/talloc-2.1.1.tar.gz
}

function download_libarchive()
{
  cd ${CUR_PATH}
  rm -rf libarchive
  #curl -o libarchive-3.5.2.tar.gz https://codeload.github.com/libarchive/libarchive/tar.gz/refs/tags/v3.5.2
  git clone https://github.com/libarchive/libarchive
  cd libarchive
  git checkout -b v3.5.2
}

function download_proot()
{
  cd ${CUR_PATH}
  curl -o proot-5.1.0.tar.gz https://codeload.github.com/proot-me/proot/tar.gz/refs/tags/v5.1.0
}

function clear_install_temp_file() {
     cd ${CUR_PATH}
     rm -rf proot-5.1.0.tar.gz
     rm -rf proot-5.1.0/
     rm -rf libarchive/
     rm -rf talloc-2.1.1.tar.gz
     rm -rf talloc-2.1.1/
     rm -rf build_proot.sh
}

ins_glibc_static

download_talloc
ins_libtalloc

download_libarchive
ins_libarchive

download_proot
ins_proot

clear_install_temp_file


