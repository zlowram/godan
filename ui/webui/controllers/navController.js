angular.module('Godan').controller('NavCtrl', ['$scope', '$location', function NavCtrl($scope, $location) {
	$scope.isActive = function (viewLocation) { 
        return viewLocation === $location.path();
    };
}]);
