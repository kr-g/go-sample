var mapp = angular.module('myapp', [ "ngRoute", "ui.bootstrap", 'ngAnimate', 'ngSanitize' ]);

mapp.config([ '$routeProvider', function($routeProvider) {
	
	$routeProvider
	.when( "/resttldsmod", {
		templateUrl : "resttlds.html"
	})
	.otherwise({
		templateUrl : "notfound.html"
	});
	
} ]);


mapp.controller('getalldata', [ '$scope', '$http', function($scope, $http, $uibModal ) {
	
	$scope.refresh = function(){
		$http.get('http://localhost:8082/resttlds').then(function(response) {
			$scope.alldata = response.data;
		});
	};
		
	$scope.remove = function(index){

		var uid = $scope.alldata[ index ].Id;
		console.log( "del " + uid );
		
		$http.delete('http://localhost:8082/resttlds/' + uid ).then(function(response) {
			$scope.refresh();
		});
	};
	
	$scope.inc = function(index){

		var uid = $scope.alldata[ index ].Id;
		$scope.alldata[ index ].Count++;
		console.log( "inc " + uid );
		$scope.obj = $scope.alldata[ index ];
		
		$http.post('http://localhost:8082/resttlds/' + uid, JSON.stringify( $scope.alldata[ index ] ) ).then(function(response) {
			$scope.refresh();
		});
	};
	
	$scope.dec = function(index){

		var uid = $scope.alldata[ index ].Id;
		if( $scope.alldata[ index ].Count == 0 ){ return; }
		
		$scope.alldata[ index ].Count--;
		console.log( "dec " + uid );
		$scope.obj = $scope.alldata[ index ];
		
		$http.post('http://localhost:8082/resttlds/' + uid, JSON.stringify( $scope.alldata[ index ] ) ).then(function(response) {
			$scope.refresh();
		});
	};
	
	$scope.save = function(obj){

		console.log( "save " , obj );
				
		$http.post('http://localhost:8082/resttlds/' + obj.Id, JSON.stringify( obj ) ).then(function(response) {
			$scope.refresh();
		});
	};
	
	$scope.refresh();	
	
	$scope.obj = { Id: "dummy", Name: "hhh"};
	
} ]);



mapp.controller('ModalDemoCtrl', function ($uibModal,$scope) {
	
	var $ctrl = this;

	$ctrl.remove = function (index ) {

	    var modalInstance = $uibModal.open({
			animation: true,
			templateUrl: 'myModalConfirm.html',
			controllerAs: '$ctrl',
			controller: function ($uibModalInstance) {
				
				var $ctrl = this;
				
				$ctrl.obj = $scope.alldata[index];
				
				$ctrl.ok = function () {
					$scope.remove(index);
					$uibModalInstance.close();
				};
				
				$ctrl.cancel = function () {
					console.log('Modal dismissed at: ' + new Date());
					$uibModalInstance.dismiss();
				};
				
			},
			
			resolve: {
				// empty
			}
	    });
	};

	$ctrl.edit = function (index ) {

		var modalInstance = $uibModal.open({
			animation: true,
			templateUrl: 'myModalEdit.html',
			controllerAs: '$ctrl',
			controller: function ($uibModalInstance) {

				var $ctrl = this;				
				$ctrl.obj = {};
				
				if( index >= 0 ){
					// make a copy of the original instead of editing the original
					
					//$ctrl.obj = Object.assign( {}, $scope.alldata[index] );
					// or directly using angular !
					$ctrl.obj = angular.copy( $scope.alldata[index] );
				}

				$ctrl.ok = function () {
					
					if( $ctrl.obj.Name == undefined || $ctrl.obj.Name.length == 0 ){
						return;
					}
					if( $ctrl.obj.Tld == undefined || $ctrl.obj.Tld.length == 0 ){
						return;
					}
					
					$scope.save($ctrl.obj);
					$uibModalInstance.close();
				};

				$ctrl.cancel = function () {
					console.log('Modal dismissed at: ' + new Date());
					$uibModalInstance.dismiss();
				};
			},
			resolve: {
				// empty
			}
		});	
	};

	$ctrl.newentry = function ( ) {
	
		var modalInstance = $uibModal.open({
			
			animation: true,
			templateUrl: 'myModalEdit.html',
			controllerAs: '$ctrl',
			controller: function ($uibModalInstance) {

				var $ctrl = this;	
				$ctrl.obj = { Id: '' };		
				$ctrl.ok = function () {
				
					if( $ctrl.obj.Name == undefined || $ctrl.obj.Name.length == 0 ){
						return;
					}
					if( $ctrl.obj.Tld == undefined || $ctrl.obj.Tld.length == 0 ){
						return;
					}
				
					$scope.save($ctrl.obj);
					$uibModalInstance.close();
				};
		
				$ctrl.cancel = function () {
					console.log('Modal dismissed at: ' + new Date());
					$uibModalInstance.dismiss();
				};
			},
			resolve: {
				// empty
			}
		});
	};
});


console.log("done");
