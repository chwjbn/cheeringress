import React, { useRef, useState } from 'react';
import type { ActionType, ProColumns } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import {
  CtlIngressActionBackendNodePage,
  CtlIngressActionBackendNodeRemove,
} from '@/services/ant-design-pro/api';
import { DeleteOutlined, ExclamationCircleOutlined, PlusOutlined } from '@ant-design/icons';
import { ShowErrorMessage, ShowSuccessMessage } from '@/services/ant-design-pro/lib';
import { Button, Modal } from 'antd';
import AddNodeForm from './AddNodeForm';

export type NodeListFormProps = {
  onModalVisible: (flag: boolean) => void;
  modalVisible: boolean;
  namespace_id?: string;
  backend_id?: string;
};

const NodeListForm: React.FC<NodeListFormProps> = (props) => {
  const [addModalVisible, setAddModalVisible] = useState<boolean>(false);
  const actionRef = useRef<ActionType>();

  const onFormVisible = (flag: boolean) => {
    return props.onModalVisible(flag);
  };

  const onAddModalVisible = (flag: boolean) => {
    if (!flag) {
      actionRef.current?.reload();
    }
    setAddModalVisible(flag);
  };

  const doRemoveDataAction = async (dataId?: string) => {
    const reqData: API.AppDataIdRequest = {};
    reqData.data_id = dataId;

    const respData = await CtlIngressActionBackendNodeRemove(reqData);

    if (respData.error_code === '0') {
      ShowSuccessMessage(respData.error_message);
      actionRef.current?.reload();
      return;
    }

    ShowErrorMessage(respData.error_message);
  };

  const doRemoveData = async (dataId?: string) => {
    Modal.confirm({
      title: '系统操作确认',
      icon: <ExclamationCircleOutlined />,
      content: '此数据删除操作不可恢复,你确定要删除此数据吗?',
      okText: '确认',
      cancelText: '取消',
      onOk: async () => {
        await doRemoveDataAction(dataId);
      },
    });
  };

  // 定义列
  const columns: ProColumns<API.AppDataIngressActionBackendNode>[] = [
    {
      title: '系统ID',
      dataIndex: 'data_id',
      hideInSearch: true,
      hideInTable: true,
    },

    {
      title: '节点名称',
      dataIndex: 'title',
      hideInSearch: true,
      align: 'center',
    },

    {
      title: '服务器地址',
      dataIndex: 'server_host',
      hideInSearch: true,
    },

    {
      title: '服务器端口',
      dataIndex: 'server_port',
      hideInSearch: true,
    },

    {
      title: '服务器权重',
      dataIndex: 'weight_score',
      hideInSearch: true,
    },

    {
      title: '操作',
      key: 'option',
      valueType: 'option',
      render: (_, record) => [
        <a
          key="act-del"
          onClick={() => {
            doRemoveData(record.data_id);
          }}
        >
          <DeleteOutlined />
          &nbsp;删除
        </a>,
      ],
      width: '150px',
    },
  ];

  return (
    <Modal
      title={'节点列表'}
      width="80%"
      visible={props.modalVisible}
      destroyOnClose={true}
      footer={null}
      onCancel={() => {
        onFormVisible(false);
      }}
      key="ds-modal-form"
    >
      <ProTable<API.AppDataIngressActionBackendNode, API.AppDataIngressActionBackendNodePageRequest>
        bordered={true}
        search={false}
        columns={columns}
        actionRef={actionRef}
        rowKey="data_id"
        key="p-form-data"
        params={{ backend_id: props.backend_id }}
        request={CtlIngressActionBackendNodePage}
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
            添加
          </Button>,
        ]}
      />

      <AddNodeForm
        onModalVisible={onAddModalVisible}
        modalVisible={addModalVisible}
        namespace_id={props.namespace_id}
        backend_id={props.backend_id}
      />
    </Modal>
  );
};

export default NodeListForm;
