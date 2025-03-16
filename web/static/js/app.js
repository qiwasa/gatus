(function() {
    "use strict";
    var e = {
        7488: function(e, t, s) {
            s.d(t, { L: function() { return ms } });
            s(7727);
            var n = s(9963),
                o = s(6252),
                a = s(3577),
                r = s.p + "img/logo.svg";

            const i = { class: "mb-2" },
                l = { class: "flex flex-wrap" },
                d = { class: "w-3/4 text-left my-auto" },
                g = { class: "text-3xl xl:text-5xl lg:text-4xl font-light" },
                u = { class: "w-1/4 flex justify-end" },
                h = ["src"],
                p = { key: 1, src: r, alt: "Gatus", class: "object-scale-down", style: { "max-width": "100px", "min-width": "50px", "min-height": "50px" } };

            function _(e, t, s, n, r, _) {
                const S = (0, o.up)("Loading"),
                    D = (0, o.up)("router-view"),
                    I = (0, o.up)("Tooltip"),
                    A = (0, o.up)("Social");

                return (0, o.wg)(), (0, o.iD)(o.HY, null, [
                    r.retrievedConfig ? ((0, o.wg)(), (0, o.iD)("div", {
                            key: 1,
                            class: (0, a.C_)([r.config && r.config.oidc && !r.config.authenticated ? "hidden" : "", "container container-xs relative mx-auto xl:rounded xl:border xl:shadow-xl xl:my-5 p-5 pb-12 xl:pb-5 text-left dark:bg-gray-800 dark:text-gray-200 dark:border-gray-500"]),
                            id: "global"
                        }, [
                            (0, o._)("div", i, [
                                (0, o._)("div", l, [
                                    (0, o._)("div", d, [
                                        (0, o._)("div", g, (0, a.zw)(_.header), 1)
                                    ]),
                                    (0, o._)("div", u, [
                                        ((0, o.wg)(), (0, o.j4)((0, o.LL)(_.link ? "a" : "div"), { href: _.link, target: "_blank", style: { width: "100px", "min-height": "100px" } }, {
                                            default: (0, o.w5)((() => [
                                                _.logo ? ((0, o.wg)(), (0, o.iD)("img", { key: 0, src: _.logo, alt: "Gatus", class: "object-scale-down", style: { "max-width": "100px", "min-width": "50px", "min-height": "50px" } }, null, 8, h)) :
                                                    ((0, o.wg)(), (0, o.iD)("img", p))
                                            ])),
                                            _: 1
                                        }, 8, ["href"]))
                                    ])
                                ])
                            ]),
                            (0, o.Wm)(D, { onShowTooltip: _.showTooltip }, null, 8, ["onShowTooltip"])
                        ], 2)) :
                        ((0, o.wg)(), (0, o.j4)(S, { key: 0, class: "h-64 w-64 px-4" })),
                    r.config && r.config.oidc && !r.config.authenticated ?
                        ((0, o.wg)(), (0, o.iD)("div", { key: 2, class: "mx-auto max-w-md pt-12" }, [
                            (0, o._)("img", { src: r, alt: "Gatus", class: "mx-auto", style: { "max-width": "160px", "min-width": "50px", "min-height": "50px" } }),
                            (0, o._)("h2", { class: "mt-4 text-center text-4xl font-extrabold text-gray-800 dark:text-gray-200" }, " Gatus "),
                            (0, o._)("div", { class: "py-7 px-4 rounded-sm sm:px-10" }, [
                                e.$route && e.$route.query.error ?
                                    ((0, o.wg)(), (0, o.iD)("div", { class: "text-red-500 text-center mb-5" }, [
                                        (0, o._)("div", { class: "text-sm" }, [
                                            "access_denied" === e.$route.query.error ? ((0, o.wg)(), (0, o.iD)("span", { class: "text-red-500" }, "You do not have access to this status page")) :
                                                ((0, o.wg)(), (0, o.iD)("span", { class: "text-red-500" }, (0, a.zw)(e.$route.query.error), 1))
                                        ])
                                    ])) :
                                    (0, o.kq)("", !0),
                                (0, o._)("div", null, [
                                    (0, o._)("a", { href: `${r.SERVER_URL}/oidc/login`, class: "max-w-lg mx-auto w-full flex justify-center py-3 px-4 border border-green-800 rounded-md shadow-lg text-sm text-white bg-green-700 bg-gradient-to-r from-green-600 to-green-700 hover:from-green-700 hover:to-green-800" }, " Login with OIDC ", 8)
                                ])
                            ])
                        ])) :
                        (0, o.kq)("", !0),
                    (0, o.Wm)(I, { result: r.tooltip.result, event: r.tooltip.event }, null, 8, ["result", "event"]),
                    (0, o.Wm)(A)
                ], 64)
            }

            var oe = {
                name: "App",
                components: { Loading: ne, Social: H, Tooltip: V },
                methods: {
                    fetchConfig() {
                        fetch(`${ms}/api/v1/config`, { credentials: "include" }).then((e => {
                            this.retrievedConfig = !0, 200 === e.status && e.json().then((e => { this.config = e }))
                        }))
                    },
                    showTooltip(e, t) {
                        this.tooltip = { result: e, event: t }
                    }
                },
                computed: {
                    logo() {
                        return window.config && window.config.logo && "{{ .UI.Logo }}" !== window.config.logo ? window.config.logo : ""
                    },
                    header() {
                        return window.config && window.config.header && "{{ .UI.Header }}" !== window.config.header ? window.config.header : "Health Status"
                    },
                    link() {
                        return window.config && window.config.link && "{{ .UI.Link }}" !== window.config.link ? window.config.link : null
                    },
                    buttons() {
                        return window.config && window.config.buttons ? window.config.buttons : []
                    }
                },
                data() {
                    return { error: "", retrievedConfig: !1, config: { oidc: !1, authenticated: !0 }, tooltip: {}, SERVER_URL: ms }
                },
                created() {
                    this.fetchConfig()
                }
            };

            const ae = (0, P.Z)(oe, [["render", _]]);
            var re = ae;

            (0, n.ri)(re).use(cs).mount("#app");
        }
    };
})();
