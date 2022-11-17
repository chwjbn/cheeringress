declare namespace API {
  type SelectItemNode = {
    label: string;
    value: any;
  };

  type DataMapNode = {
    data_id: string;
    data_name: string;
  };

  type BaseDataResponse = {
    success?: boolean;
    error_code?: string;
    error_message?: string;
  };

  type BasePageDataRequest = {
    current?: number;
    pageSize?: number;
  };

  type BasePageDataResponse = {
    total?: number;
    success?: boolean;
  };

  // 图形验证码
  type CheckCodeImageData = {
    code_id?: string;
    image_data?: string;
  };

  type AccountCheckCodeImageResponse = BaseDataResponse & {
    data?: CheckCodeImageData;
  };

  type AppData = {
    data_id?: string;
    tenant_id?: string;
    create_time?: string;
    update_time?: string;
    update_ip?: string;
    state?: string;
  };

  //账号
  type AccessTokenData = AppData & {
    token_data?: string;
    account_id?: string;
    username?: string;
    nickname?: string;
    real_name?: string;
    avatar?: string;
    role?: string;
  };

  type AccountLoginRequest = {
    image_code_id?: string;
    image_code_data?: string;
    username?: string;
    password?: string;
  };

  type AccountLoginResponse = BaseDataResponse & {
    data?: AccessTokenData;
  };

  type AppDataIdRequest = {
    data_id?: string;
  };

  type UserData = AppData & {
    username?: string;
    password?: string;
    nickname?: string;
    real_name?: string;
    avatar?: string;
    role?: string;
  };

  type UserPageRequest = AppDataIdRequest & {
    nickname?: string;
    real_name?: string;
    avatar?: string;
    password_old?: string;
    password_new?: string;
  };

  type UserPageResponse = BasePageDataResponse & {
    data?: UserData[];
  };

  type UserGetInfoResponse = BaseDataResponse & {
    data?: UserData;
  };

  type UserInfoUpdateRequest = AppDataIdRequest & {
    nickname?: string;
    real_name?: string;
  };

  type UserSecurityUpdateRequest = AppDataIdRequest & {
    password_old?: string;
    password_new?: string;
    password_new_confirm?: string;
  };

  type UserGetRoleMapResponse = BaseDataResponse & {
    data?: KvNode[];
  };

  type UserGetMapResponse = BaseDataResponse & {
    data?: DataMapNode[];
  };

  // AppDataIngressNamespace
  type AppDataIngressNamespace = AppData & {
    title?: string;
    last_ver?: string;
    last_pub_ver?: string;
  };

  type AppDataIngressNamespacePageRequest = BasePageDataRequest & {
    state?: string;
  };

  type AppDataIngressNamespacePageResponse = BasePageDataResponse & {
    data?: AppDataIngressNamespace[];
  };

  type AppDataIngressNamespaceAddRequest = AppDataIngressNamespace;

  type AppDataIngressNamespaceInfoResponse = BaseResponse & {
    data?: AppDataIngressNamespace;
  };

  type AppDataIngressNamespaceSaveRequest = AppDataIngressNamespace;

  // eslint-disable-next-line @typescript-eslint/ban-types
  type AppDataIngressNamespaceMapRequest = {};

  type AppDataIngressNamespaceMapResponse = BaseDataResponse & {
    data?: DataMapNode[];
  };

  // AppDataIngressActionBackend
  type AppDataIngressActionBackend = AppData & {
    namespace_id?: string;
    title?: string;
    balance_type?: string;
    node_count?: number;
  };

  type AppDataIngressActionBackendPageRequest = BasePageDataRequest & {
    state?: string;
  };

  type AppDataIngressActionBackendPageResponse = BasePageDataResponse & {
    data?: AppDataIngressActionBackend[];
  };

  type AppDataIngressActionBackendAddRequest = AppDataIngressActionBackend;

  type AppDataIngressActionBackendInfoResponse = BaseResponse & {
    data?: AppDataIngressActionBackend;
  };

  type AppDataIngressActionBackendSaveRequest = AppDataIngressActionBackend;

  // eslint-disable-next-line @typescript-eslint/ban-types
  type AppDataIngressActionBackendMapRequest = {};

  type AppDataIngressActionBackendMapResponse = BaseDataResponse & {
    data?: DataMapNode[];
  };

  type AppDataIngressActionBackendNode = AppData & {
    namespace_id?: string;
    backend_id?: string;
    title?: string;
    server_host?: string;
    server_port?: number;
    weight_score?: number;
  };

  type AppDataIngressActionBackendNodePageRequest = BasePageDataRequest & {
    backend_id?: string;
  };

  type AppDataIngressActionBackendNodePageResponse = BasePageDataResponse & {
    data?: AppDataIngressActionBackendNode[];
  };

  type AppDataIngressActionBackendNodeAddRequest = AppDataIngressActionBackendNode;

  //AppDataIngressActionStatic
  type AppDataIngressActionStatic = AppData & {
    namespace_id?: string;
    title?: string;
    content_type?: string;
    data_type?: string;
    data?: string;
  };

  type AppDataIngressActionStaticPageRequest = BasePageDataRequest & {
    state?: string;
  };

  type AppDataIngressActionStaticPageResponse = BasePageDataResponse & {
    data?: AppDataIngressActionStatic[];
  };

  type AppDataIngressActionStaticAddRequest = AppDataIngressActionStatic;

  type AppDataIngressActionStaticInfoResponse = BaseResponse & {
    data?: AppDataIngressActionStatic;
  };

  type AppDataIngressActionStaticSaveRequest = AppDataIngressActionStatic;

  // eslint-disable-next-line @typescript-eslint/ban-types
  type AppDataIngressActionStaticMapRequest = {};

  type AppDataIngressActionStaticMapResponse = BaseDataResponse & {
    data?: DataMapNode[];
  };

  //AppDataIngressSite
  type AppDataIngressSite = AppData & {
    namespace_id?: string;

    title?: string;
    order_no?: number;

    auth_need?: string;
    auth_user_name?: string;
    auth_password?: string;

    match_op?: string;
    match_value?: string;

    action_type?: string;
    action_value?: string;
  };

  type AppDataIngressSitePageRequest = BasePageDataRequest & {
    state?: string;
  };

  type AppDataIngressSitePageResponse = BasePageDataResponse & {
    data?: AppDataIngressSite[];
  };

  type AppDataIngressSiteAddRequest = AppDataIngressSite;

  type AppDataIngressSiteInfoResponse = BaseResponse & {
    data?: AppDataIngressSite;
  };

  type AppDataIngressSiteSaveRequest = AppDataIngressSite;

  // eslint-disable-next-line @typescript-eslint/ban-types
  type AppDataIngressSiteMapRequest = {};

  type AppDataIngressSiteMapResponse = BaseDataResponse & {
    data?: DataMapNode[];
  };

  //AppDataIngressSiteRule
  type AppDataIngressSiteRule = AppData & {
    namespace_id?: string;
    site_id?: string;

    title?: string;
    order_no?: number;

    http_method?: string;
    match_target?: string;
    match_op?: string;
    match_value?: string;

    action_type?: string;
    action_value?: string;
  };

  type AppDataIngressSiteRulePageRequest = BasePageDataRequest & {
    site_id?: string;
  };

  type AppDataIngressSiteRulePageResponse = BasePageDataResponse & {
    data?: AppDataIngressSiteRule[];
  };

  type AppDataIngressSiteRuleAddRequest = AppDataIngressSiteRule;

  type AppDataIngressSiteRuleInfoResponse = BaseResponse & {
    data?: AppDataIngressSiteRule;
  };

  type AppDataIngressSiteRuleSaveRequest = AppDataIngressSiteRule;

  // eslint-disable-next-line @typescript-eslint/ban-types
  type AppDataIngressSiteRuleMapRequest = {};

  type AppDataIngressSiteRuleMapResponse = BaseDataResponse & {
    data?: DataMapNode[];
  };

  type AliyunOSSUploadArg = {
    endpoint?: string;
    accessKeyId?: string;
    policy?: string;
    signature?: string;
    baseDir?: string;
  };

  type AliyunOSSUploadArgResponse = BaseDataResponse & {
    data?: AliyunOSSUploadArg;
  };
}
