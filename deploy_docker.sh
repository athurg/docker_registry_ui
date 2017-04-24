#!/bin/bash
#读取分支名
DEFAULT=`git symbolic-ref --short -q HEAD`
read -p "分支名(${DEFAULT}):" BRANCH
BRANCH=${BRANCH:-$DEFAULT}

COMMIT=`git show --no-patch --format=%H ${BRANCH}`

NAME="registry_auth"
IMAGENAME="docker/registry_auth:${BRANCH}"
[[ "${DOCKER_USER}" != "" ]] && IMAGENAME="${DOCKER_USER}/${IMAGENAME}"
IMAGENAME="sw.t4f.io/${IMAGENAME}"


BUILDPATH=${PWD}/build
rm -rf ${BUILDPATH} && mkdir -p ${BUILDPATH}

git archive --format=tar ${BRANCH} | tar x -C ${BUILDPATH}
[ $? -ne 0 ] && echo "失败!" && rm -rf ${BUILDPATH} && exit 1

echo 编译代码
docker run --rm -v ${BUILDPATH}:/go/ -e CGO_ENABLED=0 golang go install app
[ $? -ne 0 ] && echo "失败!" && rm -rf ${BUILDPATH} && exit 1


TMPPATH=`mktemp -d -t docker`
#构造镜像
mkdir -p ${TMPPATH}/go/bin
cp ${BUILDPATH}/bin/app ${TMPPATH}/go/bin/app
cp ${BUILDPATH}/Dockerfile.sec ${TMPPATH}/Dockerfile
rm -rf ${BUILDPATH}

echo 制作镜像
docker build -t ${IMAGENAME} --no-cache --label branch=${BRANCH} --label commit=${COMMIT} ${TMPPATH}
[ $? -ne 0 ] && echo "失败!" && rm -rf ${TMPPATH} && exit 1
rm -rf ${TMPPATH}

echo 上传镜像
docker push ${IMAGENAME}
