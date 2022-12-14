import React, { useRef, useState } from 'react';
import type { ActionType, ProColumns } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import {
  CtlIngressSiteRulePage,
  CtlIngressSiteRuleRemove,
  DataMapAppDataIngressActionType,
  DataMapHttpMethod,
  DataMapHttpTargetItem,
  DataMapRuleStringOp,
} from '@/services/ant-design-pro/api';
import {
  DeleteOutlined,
  ExclamationCircleOutlined,
  PlusOutlined,
  ProfileOutlined,
} from '@ant-design/icons';
import { ShowErrorMessage, ShowSuccessMessage } from '@/services/ant-design-pro/lib';
import { Button, Modal } from 'antd';
import AddRuleForm from './AddRuleForm';
import InfoRuleForm from './InfoRuleForm';
import { Access, useAccess } from 'umi';

export type RuleListFormProps = {
  onModalVisible: (flag: boolean) => void;
  modalVisible: boolean;
  namespace_id?: string;
  site_id?: string;
};

const RuleListForm: React.FC<RuleListFormProps> = (props) => {
  const [addModalVisible, setAddModalVisible] = useState<boolean>(false);
  const [xInfoModalVisible, setInfoModalVisible] = useState<boolean>(false);

  const actionRef = useRef<ActionType>();
  const appAccess = useAccess();

  const [currentRow, setCurrentRow] = useState<API.AppDataIngressSiteRule>();

  const onFormVisible = (flag: boolean) => {
    return props.onModalVisible(flag);
  };

  const onAddModalVisible = (flag: boolean) => {
    if (!flag) {
      actionRef.current?.reload();
    }
    setAddModalVisible(flag);
  };

  const onInfoModalVisible = (flag: boolean) => {
    if (!flag) {
      actionRef.current?.reload();
    }
    setInfoModalVisible(flag);
  };

  const doRemoveDataAction = async (dataId?: string) => {
    const reqData: API.AppDataIdRequest = {};
    reqData.data_id = dataId;

    const respData = await CtlIngressSiteRuleRemove(reqData);

    if (respData.error_code === '0') {
      ShowSuccessMessage(respData.error_message);
      actionRef.current?.reload();
      return;
    }

    ShowErrorMessage(respData.error_message);
  };

  const doRemoveData = async (dataId?: string) => {
    Modal.confirm({
      title: '??????????????????',
      icon: <ExclamationCircleOutlined />,
      content: '?????????????????????????????????,???????????????????????????????',
      okText: '??????',
      cancelText: '??????',
      onOk: async () => {
        await doRemoveDataAction(dataId);
      },
    });
  };

  // ?????????
  const columns: ProColumns<API.AppDataIngressSiteRule>[] = [
    {
      title: '??????ID',
      dataIndex: 'data_id',
      hideInSearch: true,
      hideInTable: true,
    },

    {
      title: '????????????',
      dataIndex: 'order_no',
      hideInSearch: true,
      width: '100px',
      align: 'center',
    },

    {
      title: '????????????',
      dataIndex: 'title',
      ellipsis: true,
    },

    {
      title: 'HTTP??????',
      dataIndex: 'http_method',
      valueType: 'select',
      hideInSearch: true,
      valueEnum: DataMapHttpMethod,
      width: '120px',
    },

    {
      title: '??????????????????',
      dataIndex: 'match_target',
      valueType: 'select',
      hideInSearch: true,
      valueEnum: DataMapHttpTargetItem,
      width: '150px',
    },

    {
      title: '??????????????????',
      dataIndex: 'match_op',
      valueType: 'select',
      hideInSearch: true,
      valueEnum: DataMapRuleStringOp,
      width: '150px',
    },
    {
      title: '??????????????????',
      dataIndex: 'match_value',
      hideInSearch: true,
      ellipsis: true,
    },

    {
      title: '????????????',
      dataIndex: 'action_type',
      valueType: 'select',
      hideInSearch: true,
      valueEnum: DataMapAppDataIngressActionType,
      width: '150px',
    },

    {
      title: '??????',
      key: 'option',
      valueType: 'option',
      render: (_, record) => [
        <a
          key="act-info"
          onClick={() => {
            setCurrentRow(record);
            onInfoModalVisible(true);
          }}
        >
          <ProfileOutlined />
          &nbsp;??????
        </a>,
        <Access key={'access-edit'} accessible={appAccess.canAdmin}>
          <a
            key="act-del"
            onClick={() => {
              doRemoveData(record.data_id);
            }}
          >
            <DeleteOutlined />
            &nbsp;??????
          </a>
        </Access>,
      ],
      width: '150px',
    },
  ];

  return (
    <Modal
      title={'??????????????????'}
      width="80%"
      visible={props.modalVisible}
      destroyOnClose={true}
      footer={null}
      onCancel={() => {
        onFormVisible(false);
      }}
      key="ds-modal-form"
    >
      <ProTable<API.AppDataIngressSiteRule, API.AppDataIngressSiteRulePageRequest>
        bordered={true}
        columns={columns}
        actionRef={actionRef}
        rowKey="data_id"
        key="p-form-data"
        params={{ site_id: props.site_id }}
        request={CtlIngressSiteRulePage}
        size="small"
        toolBarRender={() => [
          <Button
            type="primary"
            key="primary"
            onClick={() => {
              onAddModalVisible(true);
            }}
            icon={<PlusOutlined />}
          >
            ??????
          </Button>,
        ]}
      />

      <AddRuleForm
        onModalVisible={onAddModalVisible}
        modalVisible={addModalVisible}
        namespace_id={props.namespace_id}
        site_id={props.site_id}
      />

      <InfoRuleForm
        onModalVisible={onInfoModalVisible}
        modalVisible={xInfoModalVisible}
        namespace_id={currentRow?.namespace_id}
        id={currentRow?.data_id}
      />
    </Modal>
  );
};

export default RuleListForm;
