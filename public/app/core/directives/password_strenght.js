define([
  '../core_module',
],
function (coreModule) {
  'use strict';

  coreModule.default.directive('passwordStrength', function() {
    var template = '<div class="password-strength small" ng-if="!loginMode" ng-class="strengthClass">' +
      '<em>{{strengthText}}</em>' +
      '</div>';
    return {
      template: template,
      scope: {
        password: "=",
      },
      link: function($scope) {

        $scope.strengthClass = '';

        function passwordChanged(newValue) {
          if (!newValue) {
            $scope.strengthText = "";
            $scope.strengthClass = "hidden";
            return;
          }
          if (newValue.length < 4) {
            $scope.strengthText = "强度：弱.";
            $scope.strengthClass = "password-strength-bad";
            return;
          }
          if (newValue.length <= 8) {
            $scope.strengthText = "强度：中.";
            $scope.strengthClass = "password-strength-ok";
            return;
          }

          $scope.strengthText = "强度：强.";
          $scope.strengthClass = "password-strength-good";
        }

        $scope.$watch("password", passwordChanged);
      }
    };
  });
});
