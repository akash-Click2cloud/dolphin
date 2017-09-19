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
            $scope.builderImages = {
                'Click2Cloud Python Builder Image' :  'click2cloud/python-33-centos7',
                'Click2Cloud NodeJs Builder Image'  :   'click2cloud/nodejs-010-centos7',
                'Click2Cloud Ruby Builder Image'  :   'click2cloud/ruby-20-centos7',
                'Click2Cloud PHP Builder Image'  :   'click2cloud/php-55-centos7',
                'Click2Cloud Perl Builder Image'  :   'click2cloud/perl-516-centos7',
                'Click2Cloud .NET Core Builder Image' : 'click2cloud/aspnet-core-centos7',
                'Click2Cloud .NET 4.5 Builder Image' : 'click2cloud/aspnet-4.5-centos7'
            };
            $scope.isbisabled=false;

            //$scope.builderImage = ['centos/python-35-centos7','click2cloud/perl-516-centos7','click2cloud/nodejs-010-centos7','click2cloud/ruby-20-centos7','click2cloud/php-55-centos7'];

            $scope.buildApptocontainer = function() {
                    $('#atocBuildSpinner').show();
                    $('#notify').show();
                    $scope.isbisabled=true;
                    var BaseImage = $scope.formValues.BaseImage;
                    if (BaseImage in $scope.builderImages) {
                        BaseImage = $scope.builderImages[BaseImage];
                    }
                    var GitUrl =$scope.formValues.GitUrl;
                    var ImageName =$scope.formValues.ImageName;
                     appToContainerService.appToContainer(BaseImage, GitUrl, ImageName).then(function success(data) {
                         $scope.appToContainerlogs=data.Output;
                         console.log('controller function calling');
                         console.log($scope.appToContainerlogs);
                        Notifications.success('Image created', name);
                        //$state.reload();
                    })
                     .catch(function error(err) {
                         Notifications.error('Failure','', 'App to Container Image Build failed');
                         $scope.appToContainerlogs = err.data.err;

                     })
                     .finally(function final() {
                         $('#atocBuildSpinner').hide();
                         $('#notify').hide();
                         $scope.isbisabled=false;
                     });
                    $scope.appToContainerlogs = '';
            };
        }]);
