/**
 * Created by yangyang on 2017/11/4.
 */

var TerminalExt = {
    ExtWelcome: function () {
        this.writeln("\033[32m正在连接 请稍后...\033[0m");
    },
    ExtHandleWs: function (url) {
        var term = this;
        this._websockt = new WebSocket(url);
        this.attach(this._websockt);
        term.onerror = function (err) {
            console.log(err);
        };
        this._websockt.onclose = function (ent) {
            console.log(ent);
            term._websockt.close();
        };
    },
    ExtResize: function (cols, rows) {
        this.resize(cols, rows);
        this._websockt.send(JSON.stringify({
            "cols": cols,
            "rows": rows,
        }));
    },
    ExtClose: function () {
        this._websockt.send("exit\n");
        this._websockt.close();
        this.destroy();
        $(this.container).remove();
        this.container = undefined;//标记未定义
    },
    WsOpen: function (option) {
        var url = option.url;
        this.container = option.container;
        this.container.innerHTML = "";
        this.open(this.container);
        //term.fit();
        this.ExtWelcome();
        this.ExtHandleWs(url);
    },
    ExtSetBackGroundColor: function (color) {
        var theme = {
            foreground: '#ffffff',
            background: color,
            /*cursor: '#ffffff',
             selection: 'rgba(255, 255, 255, 0.3)',
             black: '#000000',
             red: '#e06c75',
             brightRed: '#e06c75',
             green: '#A4EFA1',
             brightGreen: '#A4EFA1',
             brightYellow: '#EDDC96',
             yellow: '#EDDC96',
             magenta: '#e39ef7',
             brightMagenta: '#e39ef7',
             cyan: '#5fcbd8',
             brightBlue: '#5fcbd8',
             brightCyan: '#5fcbd8',
             blue: '#5fcbd8',
             white: '#d0d0d0',
             brightBlack: '#808080',
             brightWhite: '#ffffff'*/
        };
        this._setTheme(theme);
    }

};
Terminal.prototype = Object.assign(Terminal.prototype, TerminalExt);
//对象合并http://www.cnblogs.com/yes-V-can/p/5631645.html

var ManyTerminal = function () {
    var that = this;
    this.connectedTerminal = [];
    (function () {
        window.onload = function () {
            that.ulContainer = document.getElementById("connect-list");
            that.termContainer = document.getElementById('terminal-container');
        }
    })();
};
ManyTerminal.prototype = {
    append: function (terminal) {
        //terminal 是Terminal对象
        this.connectedTerminal.push(terminal);
    },
    createHtml: function (terminal) {
        var ip = terminal.container.id;
        var liNode = $("<li>" + ip + "</li>");
        // var signNode = $('<span class="connect-sign"></span>');
        // var liSwitchNode = $('<span class="connect-switch">' + ip +
        //     '</span>');

        var liCloseNode = $('<span class="connect-close"><i class="fui-cross-circle"></i></span>');
        liCloseNode.appendTo($(liNode));
        liNode.appendTo(this.ulContainer);
        terminal.ul = {
            liCloseNode: liCloseNode,
            liNode: liNode
        };
    },
    bindEvent: function (terminal) {
        var that = this;
        terminal.ul.liNode.click(function () {
            if (manyTerminal.currentTerminal !== undefined) {
                if (manyTerminal.currentTerminal.container !== undefined) {
                    that.termContainer.removeChild(manyTerminal.currentTerminal.container);
                }
            }
            //移除当前节点
            that.termContainer.appendChild(terminal.container);
            //显示当前节点
            //
            that.setCurrentTerminal(terminal);
            that.signCurrentNode();
        });
        terminal.ul.liCloseNode.click(function () {
            terminal.ExtClose();
            terminal.ul.liNode.remove();
            //$(this).remove();
        });

    },
    setCurrentTerminal: function (terminal) {
        this.currentTerminal = terminal;
        this.signCurrentNode();
    },
    display: function (terminal) {
        this.append(terminal);
        this.createHtml(terminal);
        this.bindEvent(terminal);
        this.setCurrentTerminal(terminal);
    },
    signCurrentNode: function () {
        for (var i = 0; i < this.connectedTerminal.length; i++) {
            this.connectedTerminal[i].ul.liNode.removeClass("active")
        }
        this.currentTerminal.ul.liNode.addClass("active");
    }
};

var manyTerminal = new ManyTerminal();
var OpenTerminal = function (ip) {
    var containerWidth = 1571;
    var containerHeight = 815;
    //办公司电脑宽高

    var url = (location.protocol === "http:" ? "ws" : "wss") + "://" + location.host + "/host/ws";
    var cols = Math.floor((containerWidth - 30) / 9); //根据容器的宽高计算xterm的字符个数 然后请求终端
    var rows = Math.floor(containerHeight / 17);
    var getArgs = $.param({
        cols: cols,
        rows: rows,
        ip: ip
    });
    url = url + "?" + getArgs;
    var term = new Terminal({
        cursorBlink: false,  // Do not blink the terminal's cursor
        cols: cols,  // Set the terminal's width to 120 columns
        rows: rows  // Set the terminal's height to 80 rows
    });
    var oneTerminalContainer = document.createElement("div");
    //oneTerminalContainer.classList = ["one-terminal-container"];
    $(oneTerminalContainer).addClass("one-terminal-container").css({
        "width": containerWidth,
        "height": containerHeight
    });

    oneTerminalContainer.id = ip;
    var termContainer = document.getElementById('terminal-container');
    if (manyTerminal.currentTerminal !== undefined) {
        if (manyTerminal.currentTerminal.container !== undefined) {
            termContainer.removeChild(manyTerminal.currentTerminal.container);
        }
    }

    termContainer.appendChild(oneTerminalContainer);
    term.WsOpen({
        url: url,
        container: oneTerminalContainer
    });
    manyTerminal.display(term);
};

window.onunload = function () {
    //http://www.jb51.net/article/30640.htm
    //刷新浏览器关闭ws通道
    for (var i = 0; i < manyTerminal.connectedTerminal.length; i++) {
        manyTerminal.connectedTerminal[i].ExtClose();
    }
    window.oneTerm && oneTerm.ExtClose();
};
window.unbeforeunload = function () {
    //http://www.jb51.net/article/30640.htm
    //刷新浏览器关闭ws通道
    for (var i = 0; i < manyTerminal.connectedTerminal.length; i++) {
        manyTerminal.connectedTerminal[i].ExtClose();
    }
    window.oneTerm && oneTerm.ExtClose();
};
$(function () {
    var ipInput = $("[name=ip]").keydown(function (event) {
        if (event.keyCode === 13) {
            $("#open").click();
        }
    });
    $("#open").click(function () {
        var that = this;
        var ip = ipInput.val();
        if (ip === "") {
            alert("input ip");
            return
        }
        OpenTerminal(ip);
    });
    // document.getElementById("resize").onclick = function () {
    //     var cols = $("[name=cols]").val();
    //     var rows = $("[name=rows]").val();
    //     //必须转成int握草 后端为uint32结构体
    //     // resize json: cannot unmarshal string into Go struct field Terminal.cols of type uint32
    //     cols = parseInt(cols);
    //     rows = parseInt(rows);
    //     term.ExtResize(cols, rows);
    // };
    // document.getElementById("go-top").onclick = function () {
    //     manyTerminal.currentTerminal && manyTerminal.currentTerminal.scrollToTop();
    // };
    // document.getElementById("go-down").onclick = function () {
    //     manyTerminal.currentTerminal && manyTerminal.currentTerminal.scrollToBottom();
    // };

});
