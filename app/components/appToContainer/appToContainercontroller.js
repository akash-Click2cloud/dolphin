/*eslint linebreak-style: [2, "windows"]*/
angular.module('appToContainer', [])
    .controller('appToContainerController', ['$scope', '$state', 'appToContainerService','AppToContainerProvider', 'Notifications', 'Pagination',
        function ($scope, $state, appToContainerService,AppToContainerProvider, Notifications, Pagination) {
            $scope.appToContainerlogs = '';
            $scope.formValues = {
                BaseImage: '',
                GitUrl: '',
                ImageName: ''
            };
            $scope.builderImage = ['centos/python-35-centos7','openshift/perl-516-centos7','openshift/nodejs-010-centos7','openshift/ruby-20-centos7','openshift/php-55-centos7'];
            $scope.buildApptocontainer = function() {
                    var BaseImage =$scope.formValues.BaseImage;
                    var GitUrl =$scope.formValues.GitUrl;
                    var ImageName =$scope.formValues.ImageName;
                     appToContainerService.appToContainer(BaseImage, GitUrl, ImageName).then(function success(data) {
                         $scope.appToContainerlogs=data.Output;
                         console.log('controller function calling');
                         console.log($scope.appToContainerlogs);
                        Notifications.success('Image created', name);
                        //$state.reload();
                    }, function error(err) {
                        $scope.state.uploadInProgress = false;
                        $scope.state.error = err.msg;
                    }, function update(evt) {
                        if (evt.upload) {
                            $scope.state.uploadInProgress = evt.upload;
                        }
                    });
            };
        }]);
