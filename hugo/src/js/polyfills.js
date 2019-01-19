import smoothscroll from 'smoothscroll-polyfill';

smoothscroll.polyfill();
require('custom-event-polyfill');
require('intersection-observer');

// https://developer.mozilla.org/en-US/docs/Web/API/NodeList/forEach#Polyfill
if (window.NodeList && !NodeList.prototype.forEach) {
  NodeList.prototype.forEach = Array.prototype.forEach;
}

// from:https://github.com/jserz/js_piece/blob/master/DOM/ChildNode/remove()/remove().md
(function (arr) {
  arr.forEach(function (item) {
    if (item.hasOwnProperty('remove')) {
      return;
    }
    Object.defineProperty(item, 'remove', {
      configurable: true,
      enumerable: true,
      writable: true,
      value: function remove() {
        if (this.parentNode !== null)
          this.parentNode.removeChild(this);
      },
    });
  });
})([Element.prototype, CharacterData.prototype, DocumentType.prototype]);

// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/String/endsWith#Polyfill
if (!String.prototype.endsWith) {
  String.prototype.endsWith = function(search, this_len) {
    if (this_len === undefined || this_len > this.length) {
      this_len = this.length;
    }
    return this.substring(this_len - search.length, this_len) === search;
  };
}

// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/from#Polyfill
// Production steps of ECMA-262, Edition 6, 22.1.2.1
if (!Array.from) {
  Array.from = (function () {
    var toStr = Object.prototype.toString;
    var isCallable = function (fn) {
      return typeof fn === 'function' || toStr.call(fn) === '[object Function]';
    };
    var toInteger = function (value) {
      var number = Number(value);
      if (isNaN(number)) { return 0; }
      if (number === 0 || !isFinite(number)) { return number; }
      return (number > 0 ? 1 : -1) * Math.floor(Math.abs(number));
    };
    var maxSafeInteger = Math.pow(2, 53) - 1;
    var toLength = function (value) {
      var len = toInteger(value);
      return Math.min(Math.max(len, 0), maxSafeInteger);
    };

    // The length property of the from method is 1.
    return function from(arrayLike/*, mapFn, thisArg */) {
      // 1. Let C be the this value.
      var C = this;

      // 2. Let items be ToObject(arrayLike).
      var items = Object(arrayLike);

      // 3. ReturnIfAbrupt(items).
      if (arrayLike == null) {
        throw new TypeError('Array.from requires an array-like object - not null or undefined');
      }

      // 4. If mapfn is undefined, then let mapping be false.
      var mapFn = arguments.length > 1 ? arguments[1] : void undefined;
      var T;
      if (typeof mapFn !== 'undefined') {
        // 5. else
        // 5. a If IsCallable(mapfn) is false, throw a TypeError exception.
        if (!isCallable(mapFn)) {
          throw new TypeError('Array.from: when provided, the second argument must be a function');
        }

        // 5. b. If thisArg was supplied, let T be thisArg; else let T be undefined.
        if (arguments.length > 2) {
          T = arguments[2];
        }
      }

      // 10. Let lenValue be Get(items, "length").
      // 11. Let len be ToLength(lenValue).
      var len = toLength(items.length);

      // 13. If IsConstructor(C) is true, then
      // 13. a. Let A be the result of calling the [[Construct]] internal method
      // of C with an argument list containing the single item len.
      // 14. a. Else, Let A be ArrayCreate(len).
      var A = isCallable(C) ? Object(new C(len)) : new Array(len);

      // 16. Let k be 0.
      var k = 0;
      // 17. Repeat, while k < len… (also steps a - h)
      var kValue;
      while (k < len) {
        kValue = items[k];
        if (mapFn) {
          A[k] = typeof T === 'undefined' ? mapFn(kValue, k) : mapFn.call(T, kValue, k);
        } else {
          A[k] = kValue;
        }
        k += 1;
      }
      // 18. Let putStatus be Put(A, "length", len, true).
      A.length = len;
      // 20. Return A.
      return A;
    };
  }());
}

// https://developer.mozilla.org/en-US/docs/Web/API/Element/classList#Polyfill
// 1. String.prototype.trim polyfill
if (!"".trim) String.prototype.trim = function(){ return this.replace(/^[\s]+|[\s]+$/g, ''); };
(function(window){"use strict"; // prevent global namespace pollution
  function checkIfValidClassListEntry(O, V){
    if (V === "") throw new DOMException(
      "Failed to execute '" + O + "' on 'DOMTokenList': The token provided must not be empty." );
    if((wsI=V.search(wsRE))!==-1)throw new DOMException("Failed to execute '"+O+"' on 'DOMTokenList': "+
      "The token provided ('"+V[wsI]+"') contains HTML space characters, which are not valid in tokens.");
  }
// 2. Implement the barebones DOMTokenList livelyness polyfill
  if (typeof DOMTokenList !== "function") (function(window){
    var document = window.document, Object = window.Object, hasOwnProp = Object.prototype.hasOwnProperty;
    var defineProperty = Object.defineProperty, allowTokenListConstruction = 0, skipPropChange = 0;
    var Element = window.Element, wsI = 0, wsRE = /[\11\12\14\15\40]/; // WhiteSpace Regular Expression
    function DOMTokenList(){
      if (!allowTokenListConstruction) throw TypeError("Illegal constructor"); // internally let it through
    }
    DOMTokenList.prototype.toString = DOMTokenList.prototype.toLocaleString = function(){return this.value};
    DOMTokenList.prototype.add = function(){
      a: for(var v=0, argLen=arguments.length,val="",ele=this["uCL"],proto=ele[" uCLp"]; v!==argLen; ++v) {
        val = arguments[v] + "", checkIfValidClassListEntry("add", val);
        for (var i=0, Len=proto.length, resStr=val; i !== Len; ++i)
          if (this[i] === val) continue a; else resStr += " " + this[i];
        this[Len] = val, proto.length += 1, proto.value = resStr;
      }
      skipPropChange = 1, ele.className = proto.value, skipPropChange = 0;
    };
    DOMTokenList.prototype.remove = function(){
      for (var v=0, argLen=arguments.length,val="",ele=this["uCL"],proto=ele[" uCLp"]; v !== argLen; ++v) {
        val = arguments[v] + "", checkIfValidClassListEntry("remove", val);
        for (var i=0, Len=proto.length, resStr="", is=0; i !== Len; ++i)
          if(is){ this[i-1]=this[i] }else{ if(this[i] !== val){ resStr+=this[i]+" "; }else{ is=1; } }
        if (!is) continue;
        delete this[Len], proto.length -= 1, proto.value = resStr;
      }
      skipPropChange = 1, ele.className = proto.value, skipPropChange = 0;
    };
    window.DOMTokenList = DOMTokenList;
    function whenPropChanges(){
      var evt = window.event, prop = evt.propertyName;
      if ( !skipPropChange && (prop==="className" || (prop==="classList" && !defineProperty)) ) {
        var target = evt.srcElement, protoObjProto = target[" uCLp"], strval = "" + target[prop];
        var tokens=strval.trim().split(wsRE), resTokenList=target[prop==="classList"?" uCL":"classList"];
        var oldLen = protoObjProto.length;
        a: for(var cI = 0, cLen = protoObjProto.length = tokens.length, sub = 0; cI !== cLen; ++cI){
          for(var innerI=0; innerI!==cI; ++innerI) if(tokens[innerI]===tokens[cI]) {sub++; continue a;}
          resTokenList[cI-sub] = tokens[cI];
        }
        for (var i=cLen-sub; i < oldLen; ++i) delete resTokenList[i]; //remove trailing indexs
        if(prop !== "classList") return;
        skipPropChange = 1, target.classList = resTokenList, target.className = strval;
        skipPropChange = 0, resTokenList.length = tokens.length - sub;
      }
    }
    function polyfillClassList(ele){
      if (!ele || !("innerHTML" in ele)) throw TypeError("Illegal invocation");
      srcEle.detachEvent( "onpropertychange", whenPropChanges ); // prevent duplicate handler infinite loop
      allowTokenListConstruction = 1;
      try{ function protoObj(){} protoObj.prototype = new DOMTokenList(); }
      finally { allowTokenListConstruction = 0 }
      var protoObjProto = protoObj.prototype, resTokenList = new protoObj();
      a: for(var toks=ele.className.trim().split(wsRE), cI=0, cLen=toks.length, sub=0; cI !== cLen; ++cI){
        for (var innerI=0; innerI !== cI; ++innerI) if (toks[innerI] === toks[cI]) { sub++; continue a; }
        this[cI-sub] = toks[cI];
      }
      protoObjProto.length = Len-sub, protoObjProto.value = ele.className, protoObjProto[" uCL"] = ele;
      if (defineProperty) { defineProperty(ele, "classList", { // IE8 & IE9 allow defineProperty on the DOM
        enumerable:   1, get: function(){return resTokenList},
        configurable: 0, set: function(newVal){
          skipPropChange = 1, ele.className = protoObjProto.value = (newVal += ""), skipPropChange = 0;
          var toks = newVal.trim().split(wsRE), oldLen = protoObjProto.length;
          a: for(var cI = 0, cLen = protoObjProto.length = toks.length, sub = 0; cI !== cLen; ++cI){
            for(var innerI=0; innerI!==cI; ++innerI) if(toks[innerI]===toks[cI]) {sub++; continue a;}
            resTokenList[cI-sub] = toks[cI];
          }
          for (var i=cLen-sub; i < oldLen; ++i) delete resTokenList[i]; //remove trailing indexs
        }
      }); defineProperty(ele, " uCLp", { // for accessing the hidden prototype
        enumerable: 0, configurable: 0, writeable: 0, value: protoObj.prototype
      }); defineProperty(protoObjProto, " uCL", {
        enumerable: 0, configurable: 0, writeable: 0, value: ele
      }); } else { ele.classList=resTokenList, ele[" uCL"]=resTokenList, ele[" uCLp"]=protoObj.prototype; }
      srcEle.attachEvent( "onpropertychange", whenPropChanges );
    }
    try { // Much faster & cleaner version for IE8 & IE9:
      // Should work in IE8 because Element.prototype instanceof Node is true according to the specs
      window.Object.defineProperty(window.Element.prototype, "classList", {
        enumerable: 1,   get: function(val){
          if (!hasOwnProp.call(ele, "classList")) polyfillClassList(this);
          return this.classList;
        },
        configurable: 0, set: function(val){this.className = val}
      });
    } catch(e) { // Less performant fallback for older browsers (IE 6-8):
      window[" uCL"] = polyfillClassList;
      // the below code ensures polyfillClassList is applied to all current and future elements in the doc.
      document.documentElement.firstChild.appendChild(document.createElement('style')).styleSheet.cssText=(
        '_*{x-uCLp:expression(!this.hasOwnProperty("classList")&&window[" uCL"](this))}' + //  IE6
        '[class]{x-uCLp/**/:expression(!this.hasOwnProperty("classList")&&window[" uCL"](this))}' //IE7-8
      );
    }
  })();
// 3. Patch in unsupported methods in DOMTokenList
  (function(DOMTokenListProto, testClass){
    if (!DOMTokenListProto.item) DOMTokenListProto.item = function(i){
      function NullCheck(n) {return n===void 0 ? null : n} return NullCheck(this[i]);
    };
    if (!DOMTokenListProto.toggle || testClass.toggle("a",0)!==false) DOMTokenListProto.toggle=function(val){
      if (arguments.length > 1) return (this[arguments[1] ? "add" : "remove"](val), !!arguments[1]);
      var oldValue = this.value;
      return (this.remove(oldValue), oldValue === this.value && (this.add(val), true) /*|| false*/);
    };
    if (!DOMTokenListProto.replace || typeof testClass.replace("a", "b") !== "boolean")
      DOMTokenListProto.replace = function(oldToken, newToken){
        checkIfValidClassListEntry("replace", oldToken), checkIfValidClassListEntry("replace", newToken);
        var oldValue = this.value;
        return (this.remove(oldToken), this.value !== oldValue && (this.add(newToken), true));
      };
    if (!DOMTokenListProto.contains) DOMTokenListProto.contains = function(value){
      for (var i=0,Len=this.length; i !== Len; ++i) if (this[i] === value) return true;
      return false;
    };
    if (!DOMTokenListProto.forEach) DOMTokenListProto.forEach = function(f){
      if (arguments.length === 1) for (var i = 0, Len = this.length; i !== Len; ++i) f( this[i], i, this);
      else for (var i=0,Len=this.length,tArg=arguments[1]; i !== Len; ++i) f.call(tArg, this[i], i, this);
    };
    if (!DOMTokenListProto.entries) DOMTokenListProto.entries = function(){
      var nextIndex = 0, that = this;
      return {next: function() {
          return nextIndex<that.length ? {value: [nextIndex, that[nextIndex]], done: false} : {done: true};
        }};
    };
    if (!DOMTokenListProto.values) DOMTokenListProto.values = function(){
      var nextIndex = 0, that = this;
      return {next: function() {
          return nextIndex<that.length ? {value: that[nextIndex], done: false} : {done: true};
        }};
    };
    if (!DOMTokenListProto.keys) DOMTokenListProto.keys = function(){
      var nextIndex = 0, that = this;
      return {next: function() {
          return nextIndex<that.length ? {value: nextIndex, done: false} : {done: true};
        }};
    };
  })(window.DOMTokenList.prototype, window.document.createElement("div").classList);
})(window);
